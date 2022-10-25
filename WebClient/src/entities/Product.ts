export interface IProduct {
    product_id: number,
    product_name: string,
    category_id:   number,
    pieces_in_pack: number,
    material_id:   number,
    weight:       number,
    lenght:       number,
    width:        number,
    height:       number,
    description:  string,
    user_id:  number,
    ozon_sku: number,
    wildberries_sku: number,
}

export  interface ICategory {
    category_id:   number,
    category_name: string,
    parent_category_id:   number,
    active:  boolean,
    subcategories: ICategory[]
}

export  interface IMaterial {
    material_id:   number,
    material_name: string,
    active:  boolean,
}

export interface IFormCategory  {
    value: number
    name: string
    class: string
}

export interface IFormMaterial{
    value: number
    name: string
    class: string
}

export function createFormCategories(category: ICategory, spaces: number, form_product_categories: IFormCategory[]) {
    const spacer = "-";
    const formProductCategory: IFormCategory ={name: spacer.repeat(spaces) + category.category_name,
        value: category.category_id, class: "category"};
        if (spaces === 0 || spaces === 1){
            formProductCategory.class = "main-category"
        }
        form_product_categories.push(formProductCategory)
    if(category.subcategories){
        spaces++
        category.subcategories.forEach((category: any) => createFormCategories(category, spaces, form_product_categories))
    }
    return
} 

export function createFormMaterials(material: IMaterial, form_product_materials: IFormMaterial[]) {
    const formProductMaterial: IFormCategory ={name: material.material_name,
        value: material.material_id, class: "material"};
        form_product_materials.push(formProductMaterial)
    return
} 