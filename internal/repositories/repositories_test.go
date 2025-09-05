package repositories

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al inicializar mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Definimos los datos de prueba
	expectedEmail := "Ubuntu@gmail.com"
	password := "Ubuntu"

	// Creamos un hash de la contraseña como estaría en la base de datos
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"password"}).
		AddRow(string(hashedPassword)) // Usamos la contraseña hasheada

	mock.ExpectQuery(`SELECT password FROM users WHERE email = (.+)`).
		WithArgs(expectedEmail).
		WillReturnRows(rows)

	// Intentamos hacer login con la contraseña sin hashear
	err = repo.Login(expectedEmail, password)

	// Verificamos que no haya errores
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestLogin_Failure(t *testing.T) {
	// Creamos una conexión mock a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al inicializar mock DB: %v", err)
	}
	defer db.Close()

	// Creamos el repositorio con la conexión mock
	repo := NewUserRepository(db)

	// Definimos los datos esperados en la base de datos mockeada
	expectedEmail := "example@example.com"
	expectedPassword := "incorrecto"

	rows := sqlmock.NewRows([]string{"password"}).
		AddRow("password")

	mock.ExpectQuery(`SELECT password FROM users WHERE email = (.+)`).
		WithArgs(expectedEmail).
		WillReturnRows(rows)

	// Ejecutamos el caso de prueba negativo (contraseña incorrecta)
	err = repo.Login(expectedEmail, expectedPassword)

	// Verificamos que se haya producido un error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid password")

	// Verificamos que todas las expectativas se hayan cumplido
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRegister_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al inicializar mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	username := "testuser"
	email := "test@example.com"
	password := "password"

	// Usamos sqlmock.AnyArg() para el campo de contraseña
	mock.ExpectExec("INSERT INTO users \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(username, email, sqlmock.AnyArg()). // AnyArg() acepta cualquier valor para la contraseña hasheada
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Register(username, email, password)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestRegister_EmailAlreadyRegistered(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al inicializar mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	username := "existinguser"
	email := "existing@example.com"
	password := "password"

	// Configuramos el mock para que devuelva error de email duplicado
	mock.ExpectExec("INSERT INTO users \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(username, email, sqlmock.AnyArg()). // Usamos AnyArg() para la contraseña hasheada
		WillReturnError(fmt.Errorf("email already registered"))

	// Ejecutamos el registro
	err = repo.Register(username, email, password)

	// Verificamos que el error sea el esperado
	assert.EqualError(t, err, "email already registered")

	// Verificamos que todas las expectativas se cumplieron
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
