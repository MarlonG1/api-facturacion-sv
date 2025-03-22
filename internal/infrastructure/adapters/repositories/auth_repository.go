package repositories

import (
	"context"
	"errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"gorm.io/gorm"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/error"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) ports.AuthRepositoryPort {
	return &AuthRepository{db: db}
}

// GetAuthTypeByApiKey obtiene el tipo de autenticación de un usuario por su API key
func (r *AuthRepository) GetAuthTypeByApiKey(ctx context.Context, apiKey string) (string, error) {
	// 1. Obtener usuario por API key
	user, err := r.GetByBranchOfficeApiKey(ctx, apiKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errPackage.ErrUserNotFound
		}
	}

	return user.AuthType, nil
}

// GetByNIT obtiene un usuario por su NIT
func (r *AuthRepository) GetByNIT(ctx context.Context, nit string) (*user.User, error) {
	var dbUser db_models.User
	// 1. Obtener usuario por NIT
	result := r.db.WithContext(ctx).Where("nit = ? AND status = ?", nit, true).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errPackage.ErrUserNotFound
		}
		return nil, result.Error
	}

	// 2. Convertir a modelo de dominio
	user := &user.User{
		ID:             dbUser.ID,
		NIT:            dbUser.NIT,
		NRC:            dbUser.NRC,
		Status:         dbUser.Status,
		AuthType:       dbUser.AuthType,
		PasswordPri:    dbUser.PasswordPri,
		CommercialName: dbUser.CommercialName,
		Business:       dbUser.Business,
		Email:          dbUser.Email,
		YearInDTE:      dbUser.YearInDTE,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}

	return user, nil
}

// GetByBranchOfficeApiKey obtiene un usuario por su API key de sucursal
func (r *AuthRepository) GetByBranchOfficeApiKey(ctx context.Context, apiKey string) (*user.User, error) {
	var branch db_models.BranchOffice

	// 1. Obtener la sucursal por su API key
	result := r.db.WithContext(ctx).Where("api_key = ? AND is_active = ?", apiKey, true).First(&branch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errPackage.ErrBranchOfficeNotFound
		}
		return nil, result.Error
	}

	// 2. Obtener el usuario asociado
	var dbUser db_models.User
	result = r.db.WithContext(ctx).Where("id = ? AND status = ?", branch.UserID, true).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errPackage.ErrUserNotFound
		}
		return nil, result.Error
	}

	// 3. Convertir a modelo de dominio
	return &user.User{
		ID:                   dbUser.ID,
		NIT:                  dbUser.NIT,
		NRC:                  dbUser.NRC,
		Status:               dbUser.Status,
		AuthType:             dbUser.AuthType,
		PasswordPri:          dbUser.PasswordPri,
		CommercialName:       dbUser.CommercialName,
		EconomicActivity:     dbUser.EconomicActivity,
		EconomicActivityDesc: dbUser.EconomicActivityDesc,
		Phone:                dbUser.Phone,
		Business:             dbUser.Business,
		Email:                dbUser.Email,
		YearInDTE:            dbUser.YearInDTE,
		CreatedAt:            dbUser.CreatedAt,
		UpdatedAt:            dbUser.UpdatedAt,
	}, nil
}

// Create crea un usuario con sus sucursales
func (r *AuthRepository) Create(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Convertir modelo de dominio a modelo de base de datos
		dbUser := db_models.User{
			NIT:                  user.NIT,
			NRC:                  user.NRC,
			Status:               true,
			AuthType:             user.AuthType,
			PasswordPri:          user.PasswordPri,
			CommercialName:       user.CommercialName,
			EconomicActivity:     user.EconomicActivity,
			EconomicActivityDesc: user.EconomicActivityDesc,
			Business:             user.Business,
			Email:                user.Email,
			Phone:                user.Phone,
			YearInDTE:            user.YearInDTE,
		}

		// 2. Crear usuario
		if err := tx.Create(&dbUser).Error; err != nil {
			return err
		}

		// 3. Actualizar ID en el modelo de dominio
		user.ID = dbUser.ID

		// 4. Crear sucursales
		for i := range user.BranchOffices {
			dbBranch := db_models.BranchOffice{
				UserID:              dbUser.ID,
				EstablishmentCode:   user.BranchOffices[i].EstablishmentCode,
				EstablishmentCodeMH: user.BranchOffices[i].EstablishmentCodeMH,
				Email:               user.BranchOffices[i].Email,
				APIKey:              user.BranchOffices[i].APIKey,
				APISecret:           user.BranchOffices[i].APISecret,
				Phone:               user.BranchOffices[i].Phone,
				EstablishmentType:   user.BranchOffices[i].EstablishmentType,
				POSCode:             user.BranchOffices[i].POSCode,
				POSCodeMH:           user.BranchOffices[i].POSCodeMH,
				IsActive:            user.BranchOffices[i].IsActive,
			}

			if err := tx.Create(&dbBranch).Error; err != nil {
				return err
			}

			// 5. Si la sucursal tiene dirección, crearla también
			if user.BranchOffices[i].Address != nil {
				dbAddress := db_models.Address{
					BranchID:     dbBranch.ID,
					Municipality: user.BranchOffices[i].Address.Municipality,
					Department:   user.BranchOffices[i].Address.Department,
					Complement:   user.BranchOffices[i].Address.Complement,
				}

				if err := tx.Create(&dbAddress).Error; err != nil {
					return err
				}
			}

			// 6. Actualizar ID en el modelo de dominio
			user.BranchOffices[i].ID = dbBranch.ID
		}

		return nil
	})
}

// Update actualiza un usuario
func (r *AuthRepository) Update(ctx context.Context, user *user.User) error {
	// 1. Convertir modelo de dominio a modelo de base de datos
	dbUser := db_models.User{
		ID:             user.ID,
		NIT:            user.NIT,
		NRC:            user.NRC,
		Status:         user.Status,
		AuthType:       user.AuthType,
		PasswordPri:    user.PasswordPri,
		CommercialName: user.CommercialName,
		Business:       user.Business,
		Email:          user.Email,
		YearInDTE:      user.YearInDTE,
	}

	// 2. Actualizar usuario
	return r.db.WithContext(ctx).Model(&dbUser).Updates(dbUser).Error
}

// UpdateBranchOffices actualiza las sucursales de un usuario
func (r *AuthRepository) UpdateBranchOffices(ctx context.Context, userID uint, branches []user.BranchOffice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, branch := range branches {
			// 1. Comprobar que la sucursal pertenece al usuario
			var count int64
			tx.Model(&db_models.BranchOffice{}).Where("id = ? AND user_id = ?", branch.ID, userID).Count(&count)
			if count == 0 {
				return errPackage.ErrBranchDoesNotBelong
			}

			// 2. Actualizar sucursal
			dbBranch := db_models.BranchOffice{
				ID:                  branch.ID,
				EstablishmentCode:   branch.EstablishmentCode,
				EstablishmentCodeMH: branch.EstablishmentCodeMH,
				Email:               branch.Email,
				APIKey:              branch.APIKey,
				APISecret:           branch.APISecret,
				Phone:               branch.Phone,
				EstablishmentType:   branch.EstablishmentType,
				POSCode:             branch.POSCode,
				POSCodeMH:           branch.POSCodeMH,
				IsActive:            branch.IsActive,
			}

			if err := tx.Model(&dbBranch).Updates(dbBranch).Error; err != nil {
				return err
			}

			// 3. Si hay dirección, actualizarla también
			if branch.Address != nil {
				dbAddress := db_models.Address{
					BranchID:     branch.ID,
					Municipality: branch.Address.Municipality,
					Department:   branch.Address.Department,
					Complement:   branch.Address.Complement,
				}

				// 3.1 Actualizar o crear dirección
				var existingAddress db_models.Address
				result := tx.Where("branch_id = ?", branch.ID).First(&existingAddress)
				if result.Error != nil {
					if errors.Is(result.Error, gorm.ErrRecordNotFound) {
						// 3.2 Crear dirección si no existe
						if err := tx.Create(&dbAddress).Error; err != nil {
							return err
						}
					} else {
						return result.Error
					}
				} else {
					// 3.4 Actualizar dirección existente
					dbAddress.ID = existingAddress.ID
					if err := tx.Model(&dbAddress).Updates(dbAddress).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

// DeleteBranchOffice elimina una sucursal de un usuario
func (r *AuthRepository) DeleteBranchOffice(ctx context.Context, userID uint, branchID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Comprobar que la sucursal pertenece al usuario
		var count int64
		tx.Model(&db_models.BranchOffice{}).Where("id = ? AND user_id = ?", branchID, userID).Count(&count)
		if count == 0 {
			return errPackage.ErrBranchDoesNotBelong
		}

		// 2. Eliminar dirección primero (debido a la restricción de clave foránea)
		if err := tx.Where("branch_id = ?", branchID).Delete(&db_models.Address{}).Error; err != nil {
			return err
		}

		// 3. Eliminar sucursal
		return tx.Delete(&db_models.BranchOffice{}, branchID).Error
	})
}

func (r *AuthRepository) GetBranchByApiKey(ctx context.Context, apiKey string) (*user.BranchOffice, error) {
	var branch db_models.BranchOffice

	result := r.db.WithContext(ctx).Preload("Address").Where("api_key = ?", apiKey).First(&branch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errPackage.ErrBranchOfficeNotFound
		}
		return nil, result.Error
	}

	return &user.BranchOffice{
		ID:                  branch.ID,
		UserID:              branch.UserID,
		EstablishmentCode:   branch.EstablishmentCode,
		EstablishmentCodeMH: branch.EstablishmentCodeMH,
		Email:               branch.Email,
		APIKey:              branch.APIKey,
		APISecret:           branch.APISecret,
		Phone:               branch.Phone,
		EstablishmentType:   branch.EstablishmentType,
		POSCode:             branch.POSCode,
		POSCodeMH:           branch.POSCodeMH,
		IsActive:            branch.IsActive,
		Address: &user.Address{
			Municipality: branch.Address.Municipality,
			Department:   branch.Address.Department,
			Complement:   branch.Address.Complement,
		},
	}, nil
}

// GetAuthTypeByNIT obtiene el tipo de autenticación de un usuario por su NIT
func (r *AuthRepository) GetAuthTypeByNIT(ctx context.Context, nit string) (string, error) {
	user, err := r.GetByNIT(ctx, nit)
	if err != nil {
		return "", err
	}

	return user.AuthType, nil
}

// GetIssuerInfoByApiKey obtiene la información del usuario y sucursal formateada para el DTE de Hacienda
func (r *AuthRepository) GetIssuerInfoByApiKey(ctx context.Context, apiKey string) (*dte.IssuerDTE, error) {
	// 1. Obtener sucursal
	branch, err := r.GetBranchByApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	phone := branch.Phone
	email := branch.Email

	// 2. Obtener usuario
	user, err := r.GetByBranchOfficeApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	// 2.1 Si la sucursal no tiene correo, usar el del usuario
	if branch.Email == nil {
		email = &user.Email
	}

	// 2.2 Si la sucursal no tiene teléfono, usar el del usuario
	if branch.Phone == nil {
		phone = &user.Phone
	}

	// 2.3 Si la sucursal no tiene dirección, usar la de la casa matriz
	if branch.Address == nil {
		matrixBranch, err := r.GetMatrixBranch(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		branch.Address = matrixBranch.Address
	}

	// 3. Formatear información para DTE
	return &dte.IssuerDTE{
		NIT:                  user.NIT,
		NRC:                  user.NRC,
		CommercialName:       user.CommercialName,
		BusinessName:         user.Business,
		EconomicActivity:     user.EconomicActivity,
		EconomicActivityDesc: user.EconomicActivityDesc,
		EstablishmentCode:    branch.EstablishmentCode,
		EstablishmentCodeMH:  branch.EstablishmentCodeMH,
		EstablishmentType:    branch.EstablishmentType,
		POSCode:              branch.POSCode,
		POSCodeMH:            branch.POSCodeMH,
		Email:                email,
		Phone:                phone,
		Address:              branch.Address,
	}, nil
}

// GetMatrixBranch obtiene la sucursal registrada como casa matriz
func (r *AuthRepository) GetMatrixBranch(ctx context.Context, userID uint) (*user.BranchOffice, error) {
	var branch db_models.BranchOffice

	result := r.db.WithContext(ctx).
		Preload("Address").
		Where("user_id = ? AND establishment_type = ?", userID, constants.CasaMatriz).
		First(&branch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errPackage.ErrBranchOfficeNotFound
		}
		return nil, result.Error
	}

	return &user.BranchOffice{
		ID:                  branch.ID,
		UserID:              branch.UserID,
		EstablishmentCode:   branch.EstablishmentCode,
		EstablishmentCodeMH: branch.EstablishmentCodeMH,
		Email:               branch.Email,
		APIKey:              branch.APIKey,
		APISecret:           branch.APISecret,
		Phone:               branch.Phone,
		EstablishmentType:   branch.EstablishmentType,
		POSCode:             branch.POSCode,
		POSCodeMH:           branch.POSCodeMH,
		IsActive:            branch.IsActive,
		Address: &user.Address{
			Municipality: branch.Address.Municipality,
			Department:   branch.Address.Department,
			Complement:   branch.Address.Complement,
		},
	}, nil
}
