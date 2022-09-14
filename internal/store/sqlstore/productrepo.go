package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductRepo struct {
	store *Store
}

func (r *ProductRepo) Create(p *model.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO public.product (product_name, category_id, pieces_in_pack, material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING product_id",
		p.ProductName,
		p.CategoryID,
		p.PiecesInPack,
		p.MaterialID,
		p.Weight,
		p.Lenght,
		p.Weight,
		p.Height,
		p.Description,
		p.UserID,
		p.Active,
	).Scan(&p.ProductID)
}

func (r *ProductRepo) FindByUserId(userId int) ([]*model.Product, error) {

	var products []*model.Product

	rows, err := r.store.db.Query(
		"SELECT product_id, product_name, category_id, pieces_in_pack ,material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, active FROM public.Product WHERE active = true and user_id = $1",
		userId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		p := new(model.Product)
		if err := rows.Scan(
			&p.ProductID,
			&p.ProductName,
			&p.CategoryID,
			&p.PiecesInPack,
			&p.MaterialID,
			&p.Weight,
			&p.Lenght,
			&p.Width,
			&p.Height,
			&p.Description,
			&p.UserID,
			&p.Active,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepo) CreateCategory(c *model.Category) error {
	if err := c.ValidateCategory(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO public.category (category_name, parent_category_id, active) VALUES ($1, $2, $3) RETURNING category_id",
		c.CategoryName,
		NewNullInt(int64(c.ParentCategoryID)),
		c.Active,
	).Scan(&c.CategoryID)
}

func (r *ProductRepo) GetCategories() ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	rows, err := r.store.db.Query(
		"SELECT category_id, category_name, parent_category_id, active FROM public.category WHERE active = true ORDER BY category_id",
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		c := new(model.Category)
		if err := rows.Scan(
			&c.CategoryID,
			&c.CategoryName,
			&c.ParentCategoryID,
			&c.Active,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func NewNullInt(v int64) sql.NullInt64 {
	if v == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: v,
		Valid: true,
	}
}
