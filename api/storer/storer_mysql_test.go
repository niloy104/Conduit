package storer

import (
	"context"
	"fmt"
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

// Additional tests for Order and OrderItem can be added similarly.

func TestCreateOrder(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
					WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, sqlmock.AnyArg(), nil).
					WillReturnResult(sqlmock.NewResult(1, 1))

				for _, oi := range o.Items {
					mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
						WithArgs(oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()

				no, err := st.CreateOrder(context.Background(), o)
				require.NoError(t, err)
				require.Equal(t, o, no)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
					WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, sqlmock.AnyArg(), nil).
					WillReturnError(fmt.Errorf("error inserting order"))
				mock.ExpectRollback()
		
				mo, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)
				require.Nil(t, mo)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order item",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
					WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, sqlmock.AnyArg(), nil).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WithArgs(o.Items[0].Name, o.Items[0].Quantity, o.Items[0].Image, o.Items[0].Price, o.Items[0].ProductID, 1).
					WillReturnError(fmt.Errorf("error inserting order item"))

				mock.ExpectRollback()

				mo, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)
				require.Nil(t, mo)
			},
		},			
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySQLStorer(db)
			tc.test(t, st, mock)
		})
	}

}

func TestGetOrder(t *testing.T) {
	ois := []OrderItem{
		{
			ID:        1,
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
			OrderID:   1,
		},
		{
			ID:        2,
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     199.99,
			ProductID: 2,
			OrderID:   1,
		},
	}

	o := &Order{
		ID:            1,
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").WithArgs(o.ID).WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(ois[0].ID, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, ois[0].OrderID).
					AddRow(ois[1].ID, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, ois[1].OrderID)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(o.ID).WillReturnRows(oirows)

				mo, err := st.GetOrder(context.Background(), o.ID)
				require.NoError(t, err)
				require.Equal(t, o, mo)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").WithArgs(o.ID).WillReturnError(fmt.Errorf("error querying order"))

				_, err := st.GetOrder(context.Background(), o.ID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying order items",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").WithArgs(o.ID).WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(o.ID).WillReturnError(fmt.Errorf("error querying order items"))

				_, err := st.GetOrder(context.Background(), o.ID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySQLStorer(db)
			tc.test(t, st, mock)
		})
	}
}


func TestListOrders(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders").WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(1, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, 1).
					AddRow(2, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, 1)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(1).WillReturnRows(oirows)

				mo, err := st.ListOrders(context.Background())
				require.NoError(t, err)
				require.Len(t, mo, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying orders",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM orders").WillReturnError(fmt.Errorf("error querying orders"))

				_, err := st.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying order items",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders").WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error querying order items"))

				_, err := st.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySQLStorer(db)
			tc.test(t, st, mock)
		})
	}
}

func TestDeleteOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM orders WHERE id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := st.DeleteOrder(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order item",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order item"))
				mock.ExpectRollback()

				err := st.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM orders WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order"))
				mock.ExpectRollback()

				err := st.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySQLStorer(db)
			tc.test(t, st, mock)
		})
	}
}
