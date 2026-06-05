package repository

import (
	"context"
	"fmt"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository выполняет SQL-операции над таблицей users.
// Использует pgx/v5 с пулом соединений — все запросы параметризованы,
// что исключает SQL-инъекции на уровне драйвера.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository создаёт UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create вставляет нового пользователя и возвращает его ID.
func (r *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (login, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		user.Login, user.Email, user.PasswordHash, user.Role,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

// GetByID возвращает пользователя по ID.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, login, email, password_hash, role, created_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

// GetByLogin возвращает пользователя по логину.
// Используется при аутентификации: AuthService вызывает этот метод,
// затем проверяет bcrypt-хэш переданного пароля.
func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, login, email, password_hash, role, created_at
		 FROM users WHERE login = $1`, login,
	).Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return u, nil
}

// GetByEmail возвращает пользователя по email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, login, email, password_hash, role, created_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

// List возвращает всех пользователей (для панели тренера).
// pgx.CollectRows избавляет от ручного rows.Next() / rows.Scan() цикла.
func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, login, email, password_hash, role, created_at
		 FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		var u models.User
		err := row.Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
		return u, err
	})
}

// Update изменяет логин, email и роль пользователя.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET login = $1, email = $2, role = $3
		 WHERE id = $4`,
		user.Login, user.Email, user.Role, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete удаляет пользователя. ON DELETE CASCADE в БД удалит связанные записи.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
