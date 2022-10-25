CREATE TABLE IF NOT EXISTS public.Country(
    Country_ID bigserial not null primary key,
	Country_Name varchar(200) not null,
    Country_Code varchar(200) not null
);

CREATE TABLE IF NOT EXISTS public.Currency(
    Currency_ID bigserial not null primary key,
	Currency_Name varchar(200) not null,
    Currency_Code varchar(200) not null
);

CREATE TABLE IF NOT EXISTS public.Supplier(
    Supplier_ID bigserial not null primary key,
	Supplier_Name varchar(200) not null,
    Supplier_Address varchar(500),
    Supplier_Country_ID bigint references public.Country(Country_ID),
    Supplier_SWIFT varchar(200),
    Supplier_Account_Number bigint,
    User_ID bigint not null references public.users(id),
	Active boolean not null
);

CREATE TABLE IF NOT EXISTS public.SupplyOrderStatus(
    SupplyOrderStatus_ID bigserial not null primary key,
	SupplyOrderStatus_Name varchar(200) not null
);

CREATE TABLE IF NOT EXISTS public.SupplyOrder(
    SupplyOrder_ID bigserial not null primary key,
	SupplyOrder_Date timestamp not null,
    Supplier_ID bigint not null references public.Supplier(Supplier_ID),
    ShippingCost_To_Logistic decimal(28,3) not null,
    ShippingCost_By_Logistic decimal(28,3) not null,
    User_ID bigint not null references public.users(id),
    Active boolean not null
);

CREATE TABLE IF NOT EXISTS public.SupplyOrderProduct(
    SupplyOrderProduct_ID bigserial not null primary key,
    SupplyOrder_ID bigint not null references public.SupplyOrder(SupplyOrder_ID),
	Product_ID bigint not null references public.Product(Product_ID),
    UnitPrice decimal(28,3) not null,
    Quantity decimal(28,3) not null,
    Currency_ID bigint not null references public.Currency(Currency_ID),
    SupplyOrder_Description varchar(500)
);

CREATE TABLE IF NOT EXISTS public.SupplyOrderAudit(
    SupplyOrderAudit_ID bigserial not null primary key,
    SupplyOrder_ID bigint not null references public.SupplyOrder(SupplyOrder_ID),
    SupplyOrderAudit_Date timestamp not null,
    SupplyOrderStatus_ID bigint not null,
    Audit_user_id bigint not null references public.users(id),
    Active boolean not null
);

CREATE TABLE IF NOT EXISTS public.PaymentStatus(
    PaymentStatus_ID bigserial not null primary key,
	PaymentStatus_Name varchar(200) not null
);

CREATE TABLE IF NOT EXISTS public.Payment(
    Payment_ID bigserial not null primary key,
	Payment_Date timestamp not null,
    Payment_Amount decimal(28,3) not null,
    Currency_ID bigint not null references public.Currency(Currency_ID),
    Supplier_ID bigint not null references public.Supplier(Supplier_ID),
    SupplyOrder_ID bigint not null references public.SupplyOrder(SupplyOrder_ID),
    PaymentStatus_ID bigint not null references public.PaymentStatus(PaymentStatus_ID),
    User_ID bigint not null references public.users(id),
    Active boolean not null
);


CREATE TABLE IF NOT EXISTS public.PaymentAudit(
    PaymentAudit_ID bigserial not null primary key,
    Payment_ID bigint not null references public.SupplyOrder(SupplyOrder_ID),
    PaymentAudit_Date timestamp not null,
    PaymentStatus_ID bigint not null references public.PaymentStatus(PaymentStatus_ID),
    Audit_user_id bigint not null references public.users(id),
    Active boolean not null
);

CREATE TABLE IF NOT EXISTS public.SupplyOrderPayment(
    SupplyOrderPayment_ID bigserial not null primary key,
    SupplyOrder_ID bigint not null references public.SupplyOrder(SupplyOrder_ID),
	Payment_ID bigint not null references public.Payment(Payment_ID)
);