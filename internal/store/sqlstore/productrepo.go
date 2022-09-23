package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductRepo struct {
	store *Store
}

func (r *ProductRepo) Create(p *model.Product, mpiList *model.MarketPlaceItemsList) error {
	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}

	err = r.store.db.QueryRow(
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

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	for _, mpi := range mpiList.MPIList {
		mpi.ProductID = p.ProductID
		err := r.store.db.QueryRow(
			"INSERT INTO public.marketplaceitem (product_id, marketplace_id, sku, user_id, active) VALUES ($1, $2, $3, $4, $5) RETURNING marketplaceitem_id",
			mpi.ProductID,
			mpi.MarketPlaceID,
			mpi.SKU,
			mpi.UserID,
			mpi.Active,
		).Scan(&mpi.MarketPlaceItemID)

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *ProductRepo) FindByUserId(userId int) ([]*model.Product, error) {
	var products []*model.Product
	rows, err := r.store.db.Query(
		"SELECT product_id, product_name, category_id, pieces_in_pack ,material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, active"+
			", coalesce((select mpi.sku from public.marketplaceitem as mpi WHERE mpi.active = true and mpi.product_id = p.product_id and mpi.marketplace_id = 1), 0) "+
			", coalesce((select mpi.sku from public.marketplaceitem as mpi WHERE mpi.active = true and mpi.product_id = p.product_id and mpi.marketplace_id = 2), 0) "+
			"FROM public.product as p WHERE active = true and user_id = $1",
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
			&p.OzonSKU,
			&p.WildberriesSKU,
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
		"SELECT category_id, category_name, coalesce(parent_category_id,0), active FROM public.category WHERE active = true ORDER BY category_id",
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

func (r *ProductRepo) CreateMaterial(m *model.Material) error {
	if err := m.ValidateMaterial(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO public.material (material_name, active) VALUES ($1, $2) RETURNING material_id",
		m.MaterialName,
		m.Active,
	).Scan(&m.MaterialID)

	return nil
}

func (r *ProductRepo) GetMaterials() ([]*model.Material, error) {
	materials := make([]*model.Material, 0)

	rows, err := r.store.db.Query(
		"SELECT material_id, material_name, active FROM public.material where active = true ORDER BY material_name")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		m := new(model.Material)
		{
			if err := rows.Scan(
				&m.MaterialID,
				&m.MaterialName,
				&m.Active,
			); err != nil {
				if err == sql.ErrNoRows {
					return nil, store.ErrRecordNotFound
				}
				return nil, err
			}
		}
		materials = append(materials, m)
	}
	return materials, nil
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
