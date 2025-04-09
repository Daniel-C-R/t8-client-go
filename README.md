# T8 Client Go

El objetivo de este proyecto es poner el práctica las habilidades para el desarrollo de una aplicación "real" adquiridas al comienzo de la prácticas en TWave (control de versiones con Git, gestión de proyectos con Poetry, creación de tests, documentación del código, etc).

Este proyecto ha sido desarrollado como parte de la formación interna inicial durante las prácticas en empresa en TWave. Concretamente, se enmarca dentro de la formación específica para el aprendizaje del lenguaje de programación Go, recreando la [aplicación desarrollada en la formación de Python](https://github.com/Daniel-C-R/t8-client) e intentando familiarizarse con la sintaxis del lenguaje, librerías habituales, estructura típica de un proyecto, uso de linters, etc.

El objetivo de la aplicación consiste en obtener la forma de onda y el espectro de una señal almacenada en un dispositivo T8 (los equipos que desarrolla la empresa) por medio de su API REST. Posteriormente, a partir del *waveform* se calculará su espectro y se comparará con el obtenido de la API para comprobar que los cálculos son correctos.

Para ejecutar la aplicación, lo primero que hay que hacer es cargar en el shell actual dos variables den entorno con las credenciales del usuario, ejecutando comandos como sigue:

```shell
# Estos datos son sólo de ejemplo
export T8_CLIENT_USER="user"
export T8_CLIENT_PASSWORD="password"
```

A continuación, se ejecuta el programa principal, indicando como argumentos el host a realizar la petición, la máquina, el punto, el modo de procesamiento y la fecha del registro a consultar en formato ISO. Un ejemplo se muestra a continuación:

```shell
go run ./cmd/t8-client/main.go --host "https://lzfs45.mirror.twave.io/lzfs45/rest" --machine "LP_Turbine" --point "MAD31CY005" --pmode "AM1" --datetime "2019-04-11T18:25:54"
```

Una vez ejecutado el programa, en la carpeta `output` se verán unas gráficas. `waveform` muestra la forma de onda de la señal, `spectrum.png` el espectro de la señal obtenido desde la API del T8 y `fft_spectrum.png` el espectro calculado por el programa.
