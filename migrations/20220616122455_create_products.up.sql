CREATE TABLE public.Product(
    ProductID bigserial not null primary key,
	ProductName nvarchar(200) not null,
	CategoryID bigint not null,
	PiecesInPack int,
	MaterialID bigint,
	WeightGR decimal(28,3), 
	LenghtMM decimal(28,3), 
	WidthMM decimal(28,3), 
	HeightMM decimal(28,3),
	Description  nvarchar(4000),
	Avtive boolean
)

CREATE TABLE public.Category(
    CategoryID bigint not null primary key,
	CategoryName nvarchar(200) not null,
    ParentCategoryID bigint,
	Avtive boolean
)