# Tournify

I have archived this repository as I have decided to no longer support this project.

This is the web application for [Tournify.io](https://www.tournify.io) (English) and the corresponding Swedish site [Turnering.io](https://www.turnering.io).

Current features
 - Group tournament creation

Future features
 - Elimination Tournaments
 - Double Elimination Tournaments
 - An API for registered users
 - Ability to see previous tournaments
 - Admin statistics
 - UI customization to personalize tournament pages

This project is based on the [Golang base project](https://github.com/uberswe/golang-base-project) and uses the [Tournify](https://github.com/tournify/tournify) package to create tournaments.

## Docker

I recommend copying `docker-compose.yml` and creating a `docker-compose.local.yml` file where you can make changes. Then run the following command:

```bash
docker-compose -f docker-compose.yml -f docker-compose.local.yml up --build
```

## Translations

This project uses [go-i18n](https://github.com/nicksnyder/go-i18n) to handle translations. Only English and Swedish is currently supported but I would gladly add more languages if someone would like to contribute.

To update languages first run `goi18n extract` to update `active.en.toml`. Then run `goi18n merge active.*.toml` to generate `translate.*.toml` which can then be translated. Finally, run `goi18n merge active.*.toml translate.*.toml` to merge the translated files into the active files.

## Contributions

Contributions are welcome and greatly appreciated. Please note that I am not looking to add any more features to this project but I am happy to take care of bugfixes, updates and other suggestions. If you have a question or suggestion please feel free to [open an issue](https://github.com/tournify/web/issues/new). To contribute code, please fork this repository, make your changes on a separate branch and then [open a pull request](https://github.com/tournify/web/compare).

For security related issues please see my profile, [@uberswe](https://github.com/uberswe), for ways of contacting me privately. 

## License

Please see the `LICENSE` file in the project repository.
