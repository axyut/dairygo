## dairy management with htmx, golang, templ

For a friend's dairy. Has a simple UI to add daily transactions and view monthly reports.

### used

-   [htmx](https://htmx.org/)
-   [golang](https://golang.org/)
-   [templ](https://templ.guide/quick-start/installation/)
-   [Toastify](https://github.com/apvarun/toastify-js/blob/master/README.md)
-   [Air live reload](https://github.com/cosmtrek/air)

### setup

```env
MONGODB_URI=mongodb+srv://USERNAME:PASSWORD@CLUSTER/
DB_NAME=dairyDB
MONGO_LOCAL=mongodb://localhost:27017/
ENV=dev
PORT=3000
```

```bash
make install # to install the dependencies
make dev # to run the developmental server
make build_x && bin/app.out # to build the binary x-> linux, darwin, windows
kill -TERM $(lsof -ti:3000) # to kill the server
```

### extra

launch.bat

```bat
./app.exe
pause
```

### random screenshots

![homePage](./assets/home.png)
![popup](./assets/popUpPage.png)
![transactionPage](./assets/transPage.png)
![reports page](./assets/reportsPage.png)

<!--
https://www.gomponents.com/
https://github.com/maragudk/gomponents

https://go-fuego.github.io/fuego/  -> frameworks (supports templates/templ/gomponents)
https://github.com/go-fuego/fuego
-->
