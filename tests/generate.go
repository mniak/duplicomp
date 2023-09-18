package tests

//go:generate mockgen -package=tests -destination=mock_server_handler_test.go github.com/mniak/duplicomp/internal/samples ServerHandler
//go:generate mockgen -package=tests -destination=mock_comparator_test.go github.com/mniak/duplicomp Comparator
