
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="https://avatars.githubusercontent.com/u/193694759?s=200&v=4" alt="Logo" width="80" height="80">
  </a>
</div>

---

## Prerequisites
Before you begin, ensure you have the following installed:
- **Node.js** (v18 or higher)
- **npm** (v9 or higher)
- A code editor like **Visual Studio Code**

---

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/DefiFundr-Labs/defifundr_web_app.git
   ```
2. Navigate to the project directory:
   ```sh
   cd defifundr_web_app
   ```
3. Install dependencies:
   ```sh
   npm install
   ```

---

## Running the Project
1. Start the development server:
   ```sh
   npm start
   ```
2. Open your browser and navigate to:
   ```
   http://localhost:3000
   ```

---

## Branch and Commit Conventions

### Branch Naming Convention
Always create branches from the `develop` branch:
- For new features:
  ```sh
  git checkout develop
  git checkout -b feature/descriptive_feature_name
  ```
  Example: `git checkout -b feature/add_campaign_creation_form`

- For bug fixes:
  ```sh
  git checkout develop
  git checkout -b fix/descriptive_fix_name
  ```
  Example: `git checkout -b fix/responsive_navbar_issues`

### Commit Message Convention
Use descriptive commit messages:
- For features:
  ```sh
  git commit -m "feature: brief description of the feature"
  ```
  Example: `git commit -m "feature: implement campaign details page"`

- For fixes:
  ```sh
  git commit -m "fix: brief description of the fix"
  ```
  Example: `git commit -m "fix: resolve mobile layout issues"`

### Pull Request Guidelines
- Always compare your pull request against the `staging` branch, NOT `main`.
- Ensure your branch is up to date with the `develop` branch before creating a pull request.

---

## ESLint Configuration
1. Install ESLint and TypeScript ESLint:
   ```sh
   npm install eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin --save-dev
   ```
2. Install `eslint-plugin-react`:
   ```sh
   npm install eslint-plugin-react --save-dev
   ```
3. Create or update `eslint.config.js`:
   ```js
   // eslint.config.js
   import react from 'eslint-plugin-react'
   export default tseslint.config({
     settings: { react: { version: '18.3' } },
     plugins: { react },
     rules: {
       ...react.configs.recommended.rules,
       ...react.configs['jsx-runtime'].rules,
     },
   })
   ```

---

## Running Tests
1. Run the test suite:
   ```sh
   npm test
   ```
2. Check test coverage:
   ```sh
   npm run test:coverage
   ```

---

## Contributing
To contribute to the DefiFundr frontend:
1. Ensure you're on the `develop` branch.
2. Pull the latest changes:
   ```sh
   git pull origin develop
   ```
3. Create a new branch following the naming conventions.
4. Make your changes and commit with the specified commit message format.
5. Push to the branch:
   ```sh
   git push origin feature/your_feature_name
   ```
6. Create a pull request targeting the `staging` branch.

---

## Folder Structure
```
src/
├── components/       # Reusable UI components
├── pages/            # Application pages
├── hooks/            # Custom React hooks
├── utils/            # Utility functions
├── styles/           # Global styles and CSS modules
├── assets/           # Static assets (images, icons, etc.)
├── App.tsx           # Main application component
├── index.tsx         # Entry point
└── ...
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or feedback, feel free to reach out:

- **Telegram**: [DefiFundr | OD](https://t.me/+8RoT2I_nM6kwZjdk)
