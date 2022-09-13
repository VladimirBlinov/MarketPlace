INSERT INTO public.category(
	category_id, category_name, parent_category_id, active)
	VALUES 
    (1, 'Дом и сад', null, true),
    (101, 'Посуда', 1, true),
    (102, 'Столовая посуда', 101, true),
    (103, 'Блюда', 102, true),
    (104, 'Менажницы', 103, true),
    (105, 'Деревянные менажницы', 104, true);

INSERT INTO public.material(
	material_id, material_name, active)
	VALUES 
    (1, 'Дерево', true),
    (2, 'Пластик', true);

INSERT INTO public.users(
	id, email, encryptedpassword, userrole, active)
	VALUES 
    (1, 'admin@test.org', 'qazwsx', 1, true);

INSERT INTO public.marketplace(
	marketplace_id, marketplace_name, active)
	VALUES (1, 'ozon', true);