import { IProduct } from "../entities/Product";

export function serializeProduct(product: IProduct) {
    
    const data = JSON.stringify({
        "product_id": product.product_id,
        "product_name": product.product_name,
        "category_id": product.category_id,
        "pieces_in_pack": product.pieces_in_pack,
        "material_id": product.material_id,
        "weight": product.weight,
        "lenght": product.lenght,
        "width": product.width,
        "height": product.height,
        "description": product.description,
        "user_id": product.user_id,
        "wildberries_sku": product.wildberries_sku,
        "ozon_sku": product.ozon_sku,
    })
        
    return data
} 