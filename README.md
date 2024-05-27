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

<!--
"kill -TERM $(lsof -ti:3000)"

insert global data
price ()
quantity

export counter

1. cow, bufallo milk, dahi, paneer, kurauni, kulfi, nauni
1. rate input and edit
1. insert these auto with +
1. insert custom kharcha as well

export thekka
(baaki)insert custom people price

import

1. daily income/outgoing price insert
1. milk (choose from cow buff) -> from 2,3 peoples people&price insert place

note\* daily remaining transfered as other dairy items(dahi, paneer etc)

calc
total quantity imported,remaining,exported
total price sold, bought, to sell
daily calc, monthly calc
-->
