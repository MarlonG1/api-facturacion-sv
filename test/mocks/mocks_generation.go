package mocks

//go:generate mockgen -destination=./auth_manager_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/domain/auth AuthManager
//go:generate mockgen -destination=./dte_service_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/domain/ports DTEService
//go:generate mockgen -destination=./transmitter_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/application/ports BaseTransmitter
//go:generate mockgen -destination=./seq_number_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents SequentialNumberManager
//go:generate mockgen -destination=./dte_creator_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents DTEManager
//go:generate mockgen -destination=./contingency_manager_mock.go -package=mocks github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency ContingencyManager
