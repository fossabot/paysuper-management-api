PaySuper Management API
=====

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-brightgreen.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/paysuper/paysuper-management-api.svg?branch=master)](https://travis-ci.org/paysuper/paysuper-management-api)
[![codecov](https://codecov.io/gh/paysuper/paysuper-management-api/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-management-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/paysuper/paysuper-management-api)](https://goreportcard.com/report/github.com/paysuper/paysuper-management-api)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpaysuper%2Fpaysuper-management-api.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpaysuper%2Fpaysuper-management-api?ref=badge_shield)

PaySuper is a unique, simple payment toolkit designed to make developers self-reliant. It’s an open-source payment service 
with a highly customizable payment form, an intuitive API, and comprehensible, eye-catching reports.

Management API is the [dashboard](https://github.com/paysuper/paysuper-management-server) and [payment form](https://github.com/paysuper/paysuper-payment-form)
REST API backend. Do not handle any payment processing business logic - just proxy all requests to [billing server](https://github.com/paysuper/paysuper-billing-server)
micro-service.

## Getting Started

UNDONE

## Accounting currency

1. PSP Currency - used to save the amount of the payment transaction in the PSP accounting currency. PSP currency can 
be set using the environment variable named "PSP_ACCOUNTING_CURRENCY"

2. Merchant currency -  used to save the amount of the payment transaction in the merchant (projects owner) accounting 
currency. Merchant currency can be set using merchant settings in PSP control panel.

3. Payment system currency - used to save the amount of the payment transaction in the payment system (payment methods 
owner) accounting currency. Payment system currency can be set using payment system settings in PSP admin panel.

## Contributing
We feel that a welcoming community is important and we ask that you follow PaySuper's [Open Source Code of Conduct](https://github.com/paysuper/code-of-conduct/blob/master/README.md) in all interactions with the community.

PaySuper welcomes contributions from anyone and everyone. Please refer to each project's style and contribution guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

The master branch of this repository contains the latest stable release of this component.

 


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpaysuper%2Fpaysuper-management-api.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpaysuper%2Fpaysuper-management-api?ref=badge_large)