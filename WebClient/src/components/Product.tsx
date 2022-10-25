import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useParams } from "react-router-dom"
import {IProduct} from "../entities/Product";
import { deleteProduct, getProduct } from "../services/productService";

const Product = () => {
    const navigate = useNavigate();
    const {product_id} = useParams();
    const [product, setProduct] = useState<IProduct>();

    useEffect(() => {
        getProduct(product_id)
        .then(product => setProduct(product))
    }, [product_id])

    console.log(product_id)

    return(
        <div>
            {
                product? (
                    <>
                    <h2>Товар {product.product_name}</h2>
                    <Link to={`/products/edit/${product_id}`}> <button className="edit-product">Редактировать</button></Link>
                    <Link to="/products" onClick={() => deleteProduct(product_id).then(() => navigate('/products', {replace: true}))}><button className="delete-product">Удалить</button></Link>
                    <p>{product.description}</p>
                    </>
                    ):(
                        <h2>Товар не найден</h2>
                    )
            }  
        </div>
    )
}

export {Product}