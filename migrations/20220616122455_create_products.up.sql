CREATE TABLE public.Category(
    Category_ID bigint not null primary key,
	Category_Name varchar(200) not null,
    Parent_Category_ID bigint,
	Avtive boolean not null
);

CREATE TABLE public.Material(
    Material_ID bigint not null primary key,
	Material_Name varchar(200) not null,
	Avtive boolean not null
);

CREATE TABLE public.Product(
    Product_ID bigserial not null primary key,
	Product_Name varchar(200) not null,
	Category_ID bigint not null references public.Category(Category_ID),
	Pieces_In_Pack int,
	Material_ID bigint not null references public.Material(Material_ID),
	Weight_GR decimal(28,3), 
	Lenght_MM decimal(28,3), 
	Width_MM decimal(28,3), 
	Height_MM decimal(28,3),
	Product_Description  varchar(4000),
    User_ID bigint not null references public.users(id),
	Avtive boolean not null
);

CREATE TABLE public.MarketPlace(
    MarketPlace_ID bigserial not null primary key,
    MarketPlace_Name varchar(200),
	Avtive boolean not null
);

CREATE TABLE public.MarketPlaceItem(
    MarketPlaceItem_ID bigserial not null primary key,
	Product_ID bigint not null references public.Product(Product_ID),
    Item_Name varchar(200),
    MarketPlace_ID bigint not null references public.MarketPlace(MarketPlace_ID),
    Barcode bigint not null,
    User_ID bigint not null references public.users(id),
	Avtive boolean not null
);

