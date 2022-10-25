import axios from "axios";
import {IMaterial, IProduct} from "../entities/Product";
import {ICategory} from "../entities/Product";
import { serializeProduct } from "../serializers/productSerializer";

const API_URL = "http://localhost:8080/api/v1"

export const createProduct = (product: IProduct) => {
  const data = serializeProduct(product)
  console.log("before post" + data)
  const axiosInstance = axios.create({
      withCredentials: true
    })
  return axiosInstance.post(API_URL + "/private/product/product", data)
}

export const getProducts = () => {
    const axiosInstance = axios.create({
        withCredentials: true
      })
    return axiosInstance.get<IProduct[]>(API_URL + "/private/product/product", { withCredentials: true })
    .then((response) => {
        return response.data
    })
};

export const getProduct = (product_id: string | undefined) => {
    const axiosInstance = axios.create({
        withCredentials: true
      })

    return axiosInstance.get<IProduct>(API_URL + `/private/product/product/${product_id}`)
    .then((response) => {
        return response.data
    })
};

export const deleteProduct = (product_id: string | undefined) => {
  const axiosInstance = axios.create({ 
      withCredentials: true,
    })
  return axiosInstance.delete(API_URL + `/private/product/product/${product_id}`)
}

export const updateProduct = (product: IProduct) => {
  const data = serializeProduct(product)
  console.log("before post" + data)
  const axiosInstance = axios.create({
      withCredentials: true,
    })
  return axiosInstance.put(API_URL + `/private/product/product/${product.product_id}`, data)
}

export const getProductCategories = () => {
    const axiosInstance = axios.create({
        withCredentials: true
      })
    return axiosInstance.get<ICategory[]>(API_URL + "/private/product/category/get_categories", { withCredentials: true })
    .then((response) => {
        const sortCategories = (categories: ICategory[]) => {
            const sortedCategories: any = [];
            const categoriesMap: any = {};
            categories.forEach((category: ICategory) => {
              categoriesMap[category.category_id] = category;
            });
            categories.forEach((category: ICategory) => {
              const parent = categoriesMap[category.parent_category_id];
              if (parent) {
                if (!parent.subcategories) {
                  parent.subcategories = [];
                }
                parent.subcategories.push(category);
              } else {
                sortedCategories.push(category);
              }
            });
            return sortedCategories;
          }

        return sortCategories(response.data)
    })
}

export const getProductMaterials = () => {
  const axiosInstance = axios.create({
      withCredentials: true
    })
  return axiosInstance.get<IMaterial[]>(API_URL + "/private/product/material/get_materials", { withCredentials: true })
  .then((response) => {
      return response.data
  })
}