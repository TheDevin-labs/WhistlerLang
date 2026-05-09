
# Contributing to WhistlerLang 🐧🦆
 
- Cheers for considering a contribution! This guide explains how to safely and effectively contribute to WhistlerLang.


---

# How to Contribute

1. Fork the repository


```git
git clone https://github.com/CoolyDucks/WhistlerLang
```
2. Create a new branch for your feature or fix



git checkout -b feature/awesome-feature

3. Make your changes in source/, examples/, or docs/



Keep say as the only exception to strong typing.

Avoid breaking the REPL commands.

Update examples in examples/ folder if necessary.


4. Test your changes



Run the REPL or build scripts to ensure everything works:

```
./build.sh
./build/WhistlerLang-linux-amd64 (or aarch64... anything)
```

5. Commit with clear messages


````
git add .
git commit -m "Add feature: improved time module"
````
6. Push and create a Pull Request



git push origin feature/awesome-feature

Describe your changes clearly.

Include examples if your contribution affects scripts or REPL behaviour.



---

Guidelines

Maintain strong typing rules, except for say.

Do not remove or rename original authors in code or examples.

Avoid introducing bugs that break other platforms.

Keep REPL output predictable.

Document your changes if they affect usage or examples.



---

Reporting Issues

If you spot a bug or have a feature request, please open an issue here:

https://github.com/CoolyDucks/WhistlerLang/issues

Include:

Version of WhistlerLang

Steps to reproduce the bug

Expected vs actual behaviour



---

Code of Conduct

Be respectful and collaborative. Contributions are welcome from everyone, no matter their experience.


---

Thanks ever so much for helping make WhistlerLang better!


---
