package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/saurabhraut1212/students_api_golang/internal/types"
	"github.com/saurabhraut1212/students_api_golang/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) //decode data from req body at address of the student
		if errors.Is(err, io.EOF) {
			response.WriteJosn(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJosn(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) //typecast error
			response.WriteJosn(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJosn(w, http.StatusCreated, map[string]string{"success": "ok"})

	}
}
