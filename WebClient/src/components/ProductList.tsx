import { useEffect, useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import {IProduct} from "../entities/Product"
import { deleteProduct, getProducts } from "../services/productService"
import './ProductList.scss'

const ProductList = () => {
    const [products, setProducts] = useState<IProduct[]>([])
    const navigate = useNavigate()

    useEffect(() => {
        getProducts()
            .then(products => setProducts(products))
    }, [])

    return(
        <div className="product-list-page">
        <h3>Товары</h3>
        <Link to={`/products/new`}> <button className="create-product">Создать товар</button></Link>
        {products ? (
            <div className="product-list">
            <table className="product-list-table">
                <thead>
                <tr>
                    <th></th>
                    <th>Название</th>
                    <th>Ozon SKU</th>
                    <th>Wildberries SKU</th>
                </tr>
                </thead>
                <tbody>
                {products.map((product: IProduct) => (
                <tr className="products" key={product.product_id}>
                    <td className="product-first"><Link to={`/products/edit/${product.product_id}`}> <button className="edit-product">Редактировать</button></Link>
                    <Link to="/products" onClick={() => deleteProduct(product.product_id.toString()).then(() => navigate('/products', {replace: true}))}><button className="delete-product">Удалить</button></Link>
                    </td>
                    <td><Link to={`/products/${product.product_id}`}>{product.product_name}</Link></td>
                    {(product.ozon_sku !== 0)? (<td>{product.ozon_sku}</td>):(<td></td>)}
                    {(product.wildberries_sku !== 0)? (<td>{product.wildberries_sku}</td>):(<td></td>)}
                </tr>
                ))}
                </tbody>
            </table>
            </div>
        ):
        (<><p>У Вас еще нет товаров</p>
        </>)
        }
        
        <div></div>
        </div>
    )
}

export default ProductList