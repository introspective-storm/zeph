# Zeph
### **A fast, model/framework agnostic tool to automate away the boring parts of model testing**
**v0.0.1-alpha**

## Why Zeph?
If you ever get tired of rewriting:
- training/test splits (being daring with a 85/15 split today, aren't we?)
- validation boilerplate
- chunking/streaming logic for training or testing
- logic to handle `NaN` and `null` values
- the same statistical tests for the 834th time

Zeph handles it so you don't have to.

On top of that, the tagline above isn't a lie; bring whatever model/framework fills
your heart with joy (or forced to use at work), and we will validate
your models, as performant as possible.[^1][^2]
[^1]: to the best of our ability. Sorry if we have code that is 3ms slower than it needs to be.
[^2]: we mean Zeph. Unfortunately, we didn't invent sub-millisecond model training. If you do, feel free to send a PR.

## How does it work?
Remember when we said we didn't lie? Yeah, we actually did... sort of.
Zeph is model/framework agnostic because all it cares about are three things:
1) your data file (can be specified to be pre-split, or pure test data)
2) your model file (can be script/code, serialized, model weights, etc)
3) a project specific `.zeph.yaml` or `.zeph.yml` file
Don't worry, it's just a normal `.yaml`, and you can create it yourself. That extra `.zeph` there is just so Zeph can
tell its config apart from the others on your machine. More on the `yaml` later.

You give Zeph your data file, your model file, and tell it which tests to run. It then takes care of the rest,
and gives you a dashboard with the results, right in the terminal.
Results, charts, and graphs can also be exported - PNG, PDF, or Markdown, whatever fits your workflow.

## Install
```
curl -fsSL [need to add here]
```
