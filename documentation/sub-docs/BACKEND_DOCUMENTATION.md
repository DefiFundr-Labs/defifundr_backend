


<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="https://avatars.githubusercontent.com/u/193694759?s=200&v=4" alt="Logo" width="80" height="80">
  </a>
</div>


## Installation

To set up the backend, follow these steps:

```bash
git clone https://github.com/DefiFundr-Labs/defifundr_backend.git
cd deFICrowdFunding-Backend
go mod download
```

---

## Usage

### Using Makefile

```bash
make server
```

### Using Go

```bash
go run main.go
```

### Using Air (Hot Reload)

```bash
air
```

---

## Testing

### Run All Tests

```bash
make test
```

### Unit Tests

```bash
go test ./...
```

### Test Coverage

```bash
go test -v -cover ./...
```

---

## API Documentation

The API documentation is available via Swagger. You can view it here:

[![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)](https://defifundr-hyper.koyeb.app/swagger/index.html)

[View API Documentation](https://defifundr-hyper.koyeb.app/swagger/index.html)

---

## Database Documentation

The database schema and relationships are documented using DBdocs. You can view the database documentation here:

[DB Diagram](https://dbdocs.io/kolawoleoluwasegun567/DefiFundr)

---

## Smart Contract Address

The smart contract for DefiFundr is deployed on the Ethereum Sepolia testnet. You can view the contract details on Etherscan:

[![Ethereum](https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=Ethereum&logoColor=white)](https://sepolia.etherscan.io/address/0x574Bc33136180f0734fc3fa55379e9e28701395E#code)

[View Smart Contract on Etherscan](https://sepolia.etherscan.io/address/0x574Bc33136180f0734fc3fa55379e9e28701395E#code)

---

## API Endpoints

Here are the available API endpoints for the DefiFundr backend:

| Endpoint                           |       Functionality        | HTTP Method |
| ---------------------------------- | :------------------------: | :---------: |
| `/api/v1/campaigns`                |   Create a new campaign    |    POST     |
| `/api/v1/campaigns`                |     Get all campaigns      |     GET     |
| `/api/v1/campaigns/:id`            |    Get a campaign by id    |     GET     |
| `/api/v1/campaigns/owner`          |  Get a campaign by owner   |     GET     |
| `/api/v1/campaigns/donation/:id`   |   Get a campaign donors    |     GET     |
| `/api/v1/campaigns/donate`         |    Donate to a campaign    |    POST     |
| `/api/v1/campaigns/withdraw`       |  Withdraw from a campaign  |    POST     |
| `/api/v1/campaigns/myDonations`    |      Get my donations      |     GET     |
| `/api/v1/campaigns/categories`     |     Get all categories     |     GET     |
| `/api/v1/campaigns/categories/:id` | Get campaigns by category  |     GET     |
| `/api/v1/campaigns/search`         |  Search campaigns by name  |     GET     |
| `/api/v1/campaignsTypes`           |   Get all campaign types   |     GET     |
| `/api/v1/campaigns/latestCampaigns`|    Get latest campaigns    |     GET     |
| `/api/v1/user`                     |      Get user details      |     GET     |
| `/api/v1/user`                     |    Update user details     |    POST     |
| `/api/v1/userAddress`              |    Get user by address     |    POST     |
| `/api/v1/user/avatar`              |     Set profile avatar     |     GET     |
| `/api/v1/user/avatar/set`          |       Select avatar        |    POST     |
| `/api/v1/user/biometrics`          |       Set biometrics       |    POST     |
| `/api/v1/user/logout`              |        Logout user         |    POST     |
| `/api/v1/user/password/change`     |      Change password       |    POST     |
| `/api/v1/user/password`            |      Create password       |    POST     |
| `/api/v1/user/password/reset`      |       Reset password       |    POST     |
| `/api/v1/user/password/reset/verify`| Verify password reset code |    POST     |
| `/api/v1/user/verify`              |        Verify user         |    POST     |
| `/api/v1/user/verify/resend`       |  Resend verification code  |    POST     |
| `/api/v1/user/checkUsername`       |  Check if username exists  |    POST     |
| `/api/v1/user/privatekey`          |    Get user private key    |    POST     |
| `/api/v1/user/login`               |         Login user         |    POST     |
| `/api/v1/user/renewAccess`         |     Renew access token     |    POST     |
| `/api/v1/currentPrice`             |   Get current ETH price    |     GET     |


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or feedback, feel free to reach out:

- **Telegram**: [DefiFundr | OD](https://t.me/+8RoT2I_nM6kwZjdk)

