
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="https://avatars.githubusercontent.com/u/193694759?s=200&v=4" alt="Logo" width="80" height="80">
  </a>
</div>

# DefiFundr Mobile App

[![CI](https://github.com/demola234/tdd_weather/actions/workflows/cl.yml/badge.svg)](https://github.com/demola234/tdd_weather/actions/workflows/cl.yml)
[![codecov](https://codecov.io/gh/demola234/deFICrowdFunding-Mobile/graph/badge.svg?token=VHYGUKF9YA)](https://codecov.io/gh/demola234/deFICrowdFunding-Mobile)



---

## Project Code Architecture

### V.I.P.E.R Pattern Architecture
The DefiFundr mobile app follows the **V.I.P.E.R** (View, Interactor, Presenter, Entity, Router) architecture pattern to ensure a clean and scalable codebase. This pattern separates concerns and makes the app easier to maintain and test.

---

## Installation

### Running with Makefile

1. Clean the project:
   ```bash
   make clean
   ```
   This command removes the `build` folder and the `.dart_tool` folder.

2. Build the project:
   ```bash
   make build
   ```

3. Run the project:
   ```bash
   make run
   ```

4. Generate Freezed and other generated files:
   ```bash
   make gen
   ```

### Running with Flutter Commands

1. Clean the project:
   ```bash
   flutter clean
   ```

2. Build the project:
   ```bash
   flutter build apk
   ```

3. Run the project:
   ```bash
   flutter run
   ```

4. Generate Freezed and other generated files:
   ```bash
   flutter pub run build_runner build
   ```

---

## Testing

1. Run all tests:
   ```bash
   flutter test
   ```

2. Run tests with coverage:
   ```bash
   flutter test --coverage
   ```

3. Run tests using the Makefile:
   ```bash
   make test
   ```

---

## Features

- **Payroll Management**: Create and manage payrolls securely on the blockchain.
- **Invoice Management**: Generate and track invoices with ease.
- **Blockchain Integration**: Securely interact with the Ethereum blockchain.
- **User Authentication**: Secure login and authentication system.
- **Real-Time Updates**: Get real-time updates on transactions and balances.

---


## Branch Naming Convention

We follow a structured branch naming format:

```
[fix|feat|chore|refactor]-[issue-number]-[short-description]
```

### Examples:
- `feat-23-settings-screen`
- `fix-45-settings-bug`

---

## Contributing

To contribute to the DefiFundr mobile app:
1. Fork the repository.
2. Create a new branch following the naming conventions.
3. Make your changes and commit with descriptive messages.
4. Push your changes and create a pull request targeting the `develop` branch.

---


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or feedback, feel free to reach out:

- **Telegram**: [DefiFundr | OD](https://t.me/+8RoT2I_nM6kwZjdk)


