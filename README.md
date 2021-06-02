# covaccine-notifier

CoWIN Portal Vaccine availability notifier

covaccine-notifier periodically checks and sends email notifications for available slots for the next 7 days on CoWIN portal in a given area and age.

**Sample screenshot**

![email notification](./screenshot.png)

## Installation

### Install the pre-compiled binary

```
curl -sfL https://raw.githubusercontent.com/PrasadG193/covaccine-notifier/main/install.sh | sh
```

### Docker
```
docker pull ghcr.io/prasadg193/covaccine-notifier:v0.2.0
```

## Usage

covaccine-notifier can monitor vaccine availability either by pin-code or state and district names

```bash
$ ./covaccine-notifier --help
CoWIN Vaccine availability notifier India

Usage:
  covaccine-notifier [command]

Available Commands:
  email       Notify slots availability using Email
  help        Help about any command
  telegram    Notify slots availability using Telegram

Flags:
  -a, --age int            Search appointment for age (required)
  -d, --district string    Search by district name
  -o, --dose int           Dose preference - 1 or 2. Default: 0 (both)
  -f, --fee string         Fee preferences - free (or) paid. Default: No preference
  -h, --help               help for covaccine-notifier
  -i, --interval int       Interval to repeat the search. Default: (60) second
  -m, --min-capacity int   Filter by minimum vaccination capacity. Default: (1)
  -c, --pincode string     Search by pin code
  -s, --state string       Search by state name
  -v, --vaccine string     Vaccine preferences - covishield (or) covaxin. Default: No preference

Use "covaccine-notifier [command] --help" for more information about a command.
```
example 
```
$ ./covaccine-notifier email --help 
```

**Note:** Gmail password won't work for 2FA enabled accounts. Follow [this](https://support.google.com/accounts/answer/185833?p=InvalidSecondFactor&visit_id=637554658548216477-2576856839&rd=1) guide to generate app token password and use it with `--password` arg 

**Note:** For telegram bot integration with covaccine-notifier follow [this](./docs/telegram-integration.md).

## Examples

### Terminal

#### Search by State and District

```
covaccine-notifier email --state Maharashtra --district Akola --age 27  --username <email-id> --password <email-password>
```

#### Search by Pin Code

```
covaccine-notifier email --pincode 444002 --age 27  --username <email-id> --password <email-password>
```

#### Enable Telegram Notification

```
covaccine-notifier telegram --pincode 444002 --age 27 --token <telegram-token> --username <telegram-username>
```

### Docker

```
docker run --rm -ti ghcr.io/prasadg193/covaccine-notifier:v0.2.0  email --state Maharashtra --district Akola --age 27  --username <email-id> --password <email-password>
```

### Running on Kubernetes Cluster

If you are not willing to keep your terminal on all the time :smile:, you can also create a Pod on K8s cluster

```
kubectl run covaccine-notifier --image=ghcr.io/prasadg193/covaccine-notifier:v0.2.0 --command -- /covaccine-notifier email --state Maharashtra --district Akola --age 27  --username <email-id> --password <email-password>
```

## Contributing

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:
- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
