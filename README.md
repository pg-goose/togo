## Pràctica de Permisos

> [!NOTE]
> Pau Galopa Barroso - 2025-03-28

```
rserral@asoserver:/shared$ ls -Rla

drwxr-xrwx  4 rserral profe   4096 Oct 11 10:59 .
drwxr-xr-x 12 root    root    4096 Oct 11 10:59 ..
dr-xrwx--x  2 rserral profe   4096 Oct 11 11:18 d1
drwxrwsrwt  2 root    aso     4096 Oct 11 11:18 d2

./d1:
total 8
dr-xrwx--x 2 rserral profe   4096 Oct 11 11:18 .
drwxr-xrwx 4 rserral profe   4096 Oct 11 10:59 ..
-rwx--x-w- 1 root    rserral    6 Oct 11 11:19 f2
-r--rw-rw- 1 profe   rserral 3451 Oct 11 11:00 f1

./d2:
total 7
drwxrwsrwt 2 root    aso     4096 Oct 11 11:18 .
drwxr-xrwx 4 rserral profe   4096 Oct 11 10:59 ..
-r--rwxr-- 1 rserral profe      6 Oct 11 11:19 file

aso:~$ umask
022
```

---

1. Indica si funcionaria i per què la següent comanda:
> `profe@asoserver:/shared$ echo Hola > d1/f2`
Veient els permisos d’escriptura per a *altres*, podem concloure que qualsevol pot escriure l’arxiu (no serà tan important...) i, per tant, *l’operació tindrà èxit*

---

2. Indica si funcionaria i per què la següent comanda:
> `profe@asoserver:/shared$ rm d2/file`

Veient que el propietari es `rserral` i el directori té el *sticky bit*, per tant només el propietari del fitxer, el propietari del directori o root poden eliminar-lo.

---

3. Indica si funcionaria i per què la següent comanda  
> `profe@asoserver:/shared$ mv d1/f1 d1/f2`  

Com que `profe` pot modificar el directori i tots dos arxius són modificables per qualsevol, *sí que funcionaria la comanda*  

---

4. Quins permisos, propietari i grup tindria un fitxer creat amb aquesta comanda?
> `rserral@asoserver:/shared$ touch d2/fileASO`

- Veiem que `umask` retorna `022`
- Doncs `666 - 022 = 644`
	- `6` -> `110` -> `rw-`
	- `4` -> `100` -> `r--`
	- `4` -> `100` -> `r--`
	- `rw-r--r--`
- Tindra com propietari l'usuari que el crea `rserral`
- Tindra com a grup `aso` ja que els permisos de grup del directori son `rws` i `s` o *setgid* indica que s'hereti el grup.

---

5. Indica si funcionaria i per què la següent comanda:
> `aso@asoserver:/shared$ mv d2 d3`

- Veiem que el directori
	- pertany a `rserral` amb `rwx`,
	- té com a grup `profe` amb `r-x`,
	- i altres amb `rwx`.
- Veiem que `d2`
	- pertany a `root` amb `rwx`,
	- té com a grup `aso` amb `rws`,
	- i altres amb `rwt`.

Per editar el directori tenim permís perquè "altres" permet escriptura i lectura. Per llegir `d2` tenim permís, ja que el grup i "altres" tenen permisos de lectura. Per escriure a `d3` tenim els permisos que hem vist abans. Però `d2` té el *sticky bit*, que vol dir que només el propietari o root poden moure o eliminar el directori, per molt que tinguem permisos dins de `/shared`.

