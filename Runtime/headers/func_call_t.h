#ifndef __F_CALL_H__
#define __F_CALL_H__

#include "fresult_t.h"

/*
	Хранит окружение функции.
*/
struct env_t
{
	//Переменные внешних функций.
	struct l_term* params;

	//Локальные переменные. Содержимое постоянно изменяется.
	struct l_term* locals;
};

/*
	Хранит текущее поле видимости и предыдущие.
*/
struct field_view_t
{
	//Текущее поле видимости.
	struct l_term_chain_t* current;

	//Список всех полей видимости. Необходим для восстановления
	//текущего поля видимости при откатах.
	struct l_term* backups;
};

/*
	Хранит всю информацию о запросе на вызов.
*/
struct func_call_t
{
	//Имя функции.
	const char* funcName;

	//Ссылка на саму функцию.
	struct func_result_t (*funcPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView);

	//Окружение.
	struct env_t* env;

	//Содержит информацию о поле видимости.
	struct field_view_t* fieldOfView;

	//Точка входа.
	int entryPoint;

	//Указатель на след. запрос на вызов функции
	struct l_term* next;
};

#endif
