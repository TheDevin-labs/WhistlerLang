#Contributing to WhistlerLang


Thank you for your interest in contributing to WhistlerLang. To maintain the technical integrity of Devin’ Labs projects, all contributors must adhere to the following protocols.


- 1. Technical Standards
Primary Language: All core logic must be implemented in Go (Golang).
Architecture: We have transitioned to a fresh architecture. Do not attempt to integrate or reference deprecated legacy code from previous iterations.
Dependency Management: Ensure all modules are correctly managed via go.mod.

- 2. Development Workflow
Synchronisation: Before commencing any work, ensure your local environment is synchronised with the official repository using the Makefile (Option 4).
Modular Design: Code must be modular. Ensure that Lexers, Parsers, and Evaluators remain distinct and well-documented.
Build Verification: Contributions that fail to compile via the provided Makefile will be summarily rejected.
- 3. Pull Request Process
Briefly describe the changes and the rationale behind them.
Ensure that no redundant binaries or temporary build files are included in the commit.
Avoid the use of automated AI code generation if the output is not thoroughly verified and debugged against our current architecture.
- 4. Code of Conduct
Maintain a professional demeanour within the Zulip server and GitHub discussions. Technical excellence and clear communication are the pillars of this laboratory.
