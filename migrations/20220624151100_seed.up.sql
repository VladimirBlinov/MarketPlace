INSERT INTO public.category(
	category_id, category_name, parent_category_id, active)
	VALUES 
    (1000000000, 'Дом и сад', null, true),
    (1100000000, 'Посуда', 1000000000, true),
    (1110000000, 'Столовая посуда', 1100000000, true),
    (1111000000, 'Блюда', 1110000000, true),
    (1111100000, 'Менажницы', 1111000000, true),
    (1111110000, 'Деревянные менажницы', 1111100000, true);

INSERT INTO public.material(
	material_id, material_name, active)
	VALUES 
    (1, 'Дерево', true),
    (2, 'Пластик', true);

INSERT INTO public.users(
	id, email, encryptedpassword, userrole, active)
	VALUES 
    (1, 'admin@test.org', 'qazwsx', 1, true);