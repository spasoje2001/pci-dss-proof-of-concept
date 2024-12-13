# Proof of Concept: Serverska Aplikacija u Skladu sa PCI DSS Standardom

## Pregled

Ovaj projekat predstavlja **proof of concept** serverske aplikacije razvijene u **Go** programskom jeziku, sa primarnim ciljem zadovoljavanja što većeg broja zahteva **PCI DSS (Payment Card Industry Data Security Standard)**. Za bazu podataka koristi se **PostgreSQL**.

## Funkcionalnosti

### 1. Autentifikacija korisnika
- Implementirani su login i registracija.
- Omogućena je **dvofaktorska autentifikacija** (2FA).
- Korisnik može biti **user** ili **admin**, u zavisnosti od uloge.

### 2. Funkcionalnosti prema ulozi korisnika
- **User**:
  - **CreateCardHolder**: Simulira online plaćanje unosom podataka:
    - Ime vlasnika kartice (`name`)
    - Broj kartice (`card_number`)
    - Datum isteka (`expiration_date`)
    - CVV (`cvv`)  
    CVV se unosi samo za potrebe simulacije i nije pohranjen u bazi, u skladu sa PCI DSS zahtevima. Takođe se validira da li CVV ima 3 ili 4 cifre.  
- **Admin**:
  - **GetCardholders**: Prikuplja sve transakcije koje su korisnici kreirali.

### 3. Sigurnosne mere
- **Šifrovanje i dešifrovanje podataka**:
  - **Broj kartice (`card_number`)** je šifrovan koristeći **AES algoritam**, a ključ je sačuvan u `.env` fajlu koji nije deo repozitorijuma (lokalno podešen).
  - Koristi se **pgcrypto** za šifrovanje podataka u bazi.
- CVV nije pohranjen u bazi, već se koristi isključivo za validaciju prilikom simulacije.
- **Lozinke korisnika** se hešuju pre čuvanja u bazi.
- **Kontrola pristupa** zasnovana na ulozi korisnika (**RBAC**).
- Svi podaci se prenose putem **HTTPS** protokola.

### 4. Logovanje
- Koristi se **Logrus** biblioteka za logovanje.
- Logovi se smeštaju u fajl `app.log`.

### 5. CI/CD Pipeline
- Podesili smo **CI/CD pipeline** koji automatski proverava:
  - **Čistoću koda** koristeći alat **golangci-lint** za proveru kvaliteta koda i pridržavanje Go standarda.
  - **Bezbednosne ranjivosti** koristeći alat **gosec** za identifikaciju sigurnosnih problema u kodu.
  - **Testove za autentifikaciju i autorizaciju** kako bi se osigurala ispravnost implementacije.
  - **Obim pokrivenosti testovima** i generisanje coverage izveštaja putem `go test -coverprofile`.


## Tehnologije

- **Go** (programski jezik)
- **PostgreSQL** (baza podataka)
- **Logrus** (za logovanje)
- **pgcrypto** (za šifrovanje podataka u bazi)
