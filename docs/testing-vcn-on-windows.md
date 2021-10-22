# cas - Testing cas on windows

In order to stress test cas on windows we used [machma](https://github.com/fd0/machma)

```bash
   dir /s /b *.* | machma.exe -p 4 -- cas.exe --api-key=my-cnil-api-key n {}
```
