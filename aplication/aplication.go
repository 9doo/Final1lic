package application

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/9doo/Final1lic/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
	logger *log.Logger
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
		logger: log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (a *Application) Run() error {
	for {
		a.logger.Println("Input expression:")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			a.logger.Println("Failed to read expression from console")
			continue
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			a.logger.Println("Application was successfully closed")
			return nil
		}
		result, err := calc.Calc(text)
		if err != nil {
			a.logger.Printf("%s calculation failed with error: %v", text, err)
		} else {
			a.logger.Printf("%s = %f", text, result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stdout, "[HTTP] ", log.Ldate|log.Ltime|log.Lshortfile)

	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		logger.Printf("Bad Request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calc.Calc(request.Expression)
	///w.Header().Set("Content-Type", "application/json")
	if err != nil {
		var status int
		var errMsg string
		if errors.Is(err, calc.ErrInvalidExpression) {
			status = http.StatusUnprocessableEntity
			errMsg = calc.ErrInvalidExpression.Error()
		} else if errors.Is(err, calc.ErrDivisionByZero) {
			status = http.StatusUnprocessableEntity
			errMsg = calc.ErrDivisionByZero.Error()
		} else if errors.Is(err, calc.ErrEmptyExpression) {
			status = http.StatusUnprocessableEntity
			errMsg = calc.ErrEmptyExpression.Error()
		} else {
			status = http.StatusInternalServerError
			errMsg = "unknown error"
		}
		logger.Printf("Error: %s, Status: %d, Message: %s", request.Expression, status, err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(Response{Error: errMsg})
	} else {
		logger.Printf("Successful calculation: %s = %f", request.Expression, result)
		json.NewEncoder(w).Encode(Response{Result: fmt.Sprintf("%f", result)})
	}
}

func (a *Application) RunServer() error {
	a.logger.Println("Starting server on port " + a.config.Addr)
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
