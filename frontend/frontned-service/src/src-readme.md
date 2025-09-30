# Project Structure

This document outlines the hierarchical vision for the `src` directory

---

## Directory Breakdown

- `assets/`: This directory holds all static visual assets for the application.
  - `css/`: Global stylesheets, variables, and fonts.
  - `svg/`: SVG icons and other vector-based images.

- `components/`: Contains reusable Vue components that are not full pages.

- `pages/`: Holds the primary view components, which represent the final web pages of the application. Each file here typically corresponds to a specific route defined in the router.

- `router/`: Contains all Vue Router configuration files. This is where we define the application's URL routes and map them to their corresponding page components.

- `models/` & `types/`: These directories store all TypeScript definitions.
  - `models/`: Specifically for interface definitions. (Main focus store data structures from API calls)
  - `types/`: For storing general TypeScript types, enums, and utility types used throughout the application.

- `App.vue`: The root Vue component of the application. It serves as the main layout or "shell" where page components are rendered by the router.

- `main.ts`: The entry point of the application. This file is responsible for initializing the Vue instance and integrating plugins like the router.
