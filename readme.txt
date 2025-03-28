go mod tidy

clear; go build -o app . &&\
docker build -t go-redmine-ish:latest -f DockerfileShort . &&\
docker tag go-redmine-ish localhost:32000/go-redmine-ish-golang-app:latest &&\
docker push localhost:32000/go-redmine-ish-golang-app:latest &&\
microk8s kubectl rollout restart deploy go-redmine-ish-golang-app -n go-redmine-ish

-----
Resumen de dependencias
    Independientes: users, roles, trackers, custom_fields, settings.

    Dependen de otras tablas:

        user_roles → users y roles.

        issues → trackers, projects (opcional), users (opcional).

        comments → issues y users.

        attachments → issues, comments y users.

        custom_field_values → custom_fields.
-----

posible reforma de dominios
dominio (id,nombre,descripcion)
project -> dominio
tracker -> dominio
roles -> dominio
users -> dominio


-------------
swagger

go install github.com/swaggo/swag/cmd/swag@latest

swag init

CUIDADO cuando se pasa swag init
en docs/docs.go
docs/docs.go:977:2: unknown field LeftDelim in struct literal of type "github.com/swaggo/swag".Spec
docs/docs.go:978:2: unknown field RightDelim in struct literal of type "github.com/swaggo/swag".Spec
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
-->
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	//LeftDelim:        "{{",
	//RightDelim:       "}}",
antes de compilar