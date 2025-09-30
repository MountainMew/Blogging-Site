# ODJ information

Welcome to your ODJ Repository. In this file we give you some advices and hints that helps you in working with the ODJ
tools.

We will not update this file afterwards, and you are free to override it. But for some new-joiners it would be nice if
you have a place to keep this information in your custom readme.md file.

## Secure coding advice

To improve the product security we advise you to install further plugins to your IDE.

- [SonarQube integration](https://confluence.schwarz/x/Gj-_AQ)
- [Snyk integration](https://confluence.schwarz/x/c4l2Cw)

## Pipeline templates

Want to know more about the pipelines, templating and how to use them? Checkout
the [pipeline template documentation](https://dev.azure.com/schwarzit/schwarzit.odj-pipeline-templates/_git/odj-deliver-templates?path=/README.md&version=GBmaster&_a=preview)

## API linting support for IDEs

We support the defined specs on how the open-api spec files should be written. These conventions and specifications can
be
validated via the IDE plugins for

- [Visual Studio code](https://github.com/stoplightio/vscode-spectral)
- [IntelliJ](https://github.com/SchwarzIT/spectral-intellij-plugin)

The [Linter Rules](https://github.com/SchwarzIT/api-linter-rules) are documented in GitHub.

Pipeline validates your open-api file as well.
Checkout [pipeline template documentation](https://dev.azure.com/schwarzit/schwarzit.odj-pipeline-templates/_git/odj-deliver-templates?path=/templates/technologies/common/openapi-linting/README.md&_a=preview)

## Auto generated files in ODJ Repository
ODJ is generating some files in your repository. Do not move or delete these files as this will break your ODJ provisioning.

- **.artifactignore** is a text file that controls which files are uploaded when you publish a Pipeline Artifact.
  It is typically checked into your version control repository and the syntax is similar to that of .gitignore.
  Using the .artifactignore file can help reduce your pipeline execution time by avoiding copying files into your
  staging directory before publishing your artifacts. 
  In ODJ ppipeline templates it is used to control which files are moved from the "Build sourcecode" job to the "Build container" job.


- **Dockerfile** accommodates instructions for Docker to read and automatically build images. A Dockerfile is a text
  document that contains all the commands a user could call on the command line to assemble an image


- **odj-azure-pipeline.yml** is s Pipeline file for the ODJ created pipeline. 
  ODJ Deliver provides pipeline templates that build, test, check and deploy your sourcecode. 
  Depending on the selected runtime environment, the correct mechanism for the deployment model is used.


- **sonar-project.properties** is a file to accommodate configuration properties which should be set in your project configuration and applied when a scan is run. 
  ODJ repositories contain a sonar-project.properties file that includes key-value pairs used to set the configuration values. 
  This file will only be generated if required for the selected technology.


- **spectral.config.yml** contains Schwarz API spectral custom rules to be used with [Spectral Linter](https://github.com/stoplightio/spectral) toolset. 
  It is in use to validate and lint openapi.json file. 
  Supported API Types are Product API, Backend For Frontend and Legacy API rules sets.