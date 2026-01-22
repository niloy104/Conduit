package storer

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	fn(db, mock)
}

func TestCreateProduct(t *testing.T) {
	product := &Product{
		Name:         "test Product",
		Image:        "test.jpg",
		Category:     "test Category",
		Description:  "this is a test product",
		Rating:       5,
		NumReviews:   10,
		Price:        99.99,
		CountInStock: 50,
		CreatedAt:    time.Now(),
	}

	tcs := []struct {
		name string
		test func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)
	}{
		{
			name: "sucess",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, sqlmock.AnyArg(), nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), product)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "insert error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, sqlmock.AnyArg(), nil).
					WillReturnError(sqlmock.ErrCancelled)
				cp, err := st.CreateProduct(context.Background(), product)
				require.Error(t, err)
				require.Nil(t, cp)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "last insert id error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, sqlmock.AnyArg(), nil).
					WillReturnResult(sqlmock.NewErrorResult(sqlmock.ErrCancelled))
				cp, err := st.CreateProduct(context.Background(), product)
				require.Error(t, err)
				require.Nil(t, cp)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}

}

func TestGetProduct(t *testing.T) {
	product := &Product{
		ID:           1,
		Name:         "test Product",
		Image:        "test.jpg",
		Category:     "test Category",
		Description:  "this is a test product",
		Rating:       5,
		NumReviews:   10,
		Price:        99.99,
		CountInStock: 50,
		CreatedAt:    time.Now(),
	}

	tcs := []struct {
		name string
		test func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(product.ID, product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, product.CreatedAt, nil)

				mock.ExpectQuery("SELECT * FROM products WHERE id=?").
					WithArgs(product.ID).
					WillReturnRows(rows)

				p, err := st.GetProduct(context.Background(), product.ID)
				require.NoError(t, err)
				require.Equal(t, product, p)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "get error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").
					WithArgs(product.ID).
					WillReturnError(sqlmock.ErrCancelled)

				p, err := st.GetProduct(context.Background(), product.ID)
				require.Error(t, err)
				require.Nil(t, p)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestListProducts(t *testing.T) {
	products := []*Product{
		{
			ID:           1,
			Name:         "test Product 1",
			Image:        "test1.jpg",
			Category:     "test Category",
			Description:  "this is a test product 1",
			Rating:       5,
			NumReviews:   10,
			Price:        99.99,
			CountInStock: 50,
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			Name:         "test Product 2",
			Image:        "test2.jpg",
			Category:     "test Category",
			Description:  "this is a test product 2",
			Rating:       4,
			NumReviews:   20,
			Price:        79.99,
			CountInStock: 30,
			CreatedAt:    time.Now(),
		},
	}

	tcs := []struct {
		name string
		test func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"})
				for _, p := range products {
					rows.AddRow(p.ID, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, nil)
				}

				mock.ExpectQuery("SELECT * FROM products").
					WillReturnRows(rows)

				ps, err := st.ListProducts(context.Background())
				require.NoError(t, err)
				require.Equal(t, products, ps)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "list error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products").
					WillReturnError(sqlmock.ErrCancelled)

				ps, err := st.ListProducts(context.Background())
				require.Error(t, err)
				require.Nil(t, ps)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	now := time.Now()
	product := &Product{
		ID:           1,
		Name:         "updated Product",
		Image:        "updated.jpg",
		Category:     "updated Category",
		Description:  "this is an updated product",
		Rating:       4,
		NumReviews:   15,
		Price:        89.99,
		CountInStock: 40,
		UpdatedAt:    &now,
	}

	tcs := []struct {
		name string
		test func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").
					WithArgs(product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, sqlmock.AnyArg(), product.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				p, err := st.UpdateProduct(context.Background(), product)
				require.NoError(t, err)
				require.Equal(t, product, p)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "update error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").
					WithArgs(product.Name, product.Image, product.Category, product.Description, product.Rating, product.NumReviews, product.Price, product.CountInStock, sqlmock.AnyArg(), product.ID).
					WillReturnError(sqlmock.ErrCancelled)

				p, err := st.UpdateProduct(context.Background(), product)
				require.Error(t, err)
				require.Nil(t, p)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))

				err := st.DeleteProduct(context.Background(), int64(1))
				require.NoError(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "delete error",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").
					WithArgs(int64(1)).
					WillReturnError(sqlmock.ErrCancelled)

				err := st.DeleteProduct(context.Background(), int64(1))
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}
