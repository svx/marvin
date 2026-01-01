# Directories

## `/cmd`

Main applications for this project.

## `/internal`

You can optionally add a bit of extra structure to your internal packages to separate your shared and non-shared internal code.
Your actual application code can go in the `/internal/app` directory (e.g., `/internal/app/marvin`) and the code shared by those apps in the `/internal/pkg` directory (e.g., `/internal/pkg/myprivlib`).

## `/tools`

Supporting tools for this project. Note that these tools can import code from the `/internal` directory.
