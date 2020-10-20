# concourse-tasks

A versioned catalogue of re-usable and tested Concourse tasks.

## Usage

Include the following resource in the pipeline where you wish to consume a task:

```yaml
- name: airfix
  type: git
  source:
    uri: https://github.com/EngineerBetter/concourse-tasks.git
    tag_filter: 0.0.16
```

_Always_ use a specific tag, and do not depend on the `main` branch.

To use a task, ensure that you `get: concourse-tasks` earlier in your plan. Remember to perform any input/output mapping from the generic names _inside_ the task (on the left of the colon), to your specific names _outside_ the task (on the right of the colon).

```yaml
jobs:
  name: do-the-thing
  plan:
  - get: concourse-tasks
    # ...
  - task: tarball-files
    file: concourse-tasks/tar/task.yml
    input_mapping: { input: your-directory }
    output_mapping: { output: name-you-want }
    params:
      INCLUDE: file1 file2
      TARBALL_NAME: my-tarball
```

## Why?

Concourse is trusted to build production software and deploy production infrastructure. Whilst many tasks are seemingly simple, scripting languages like Bash make it very easy for seasoned programmers to make mistakes.

Given:

* the production workloads that Concourse is entrusted with in enterprises, governments and militaries, **all tasks should be tested**.
* how long infrastructure-deploying pipelines can take to run, it is in the interests of productivity to ensure that tasks don't fail.
* the large number of pipelines that require similar tasks, we should **leverage scale by sharing common tasks**.

## Contributing

Contributions are very welcome!

1. Create your task in a descriptive directory name. No cute aeronautical puns, please.
1. Place your task config as a file called `task.yml` in your task's directory
1. Consider using the smallest image possible for your task.
1. When invoking a shell, consider [in-lining your script as an argument to the shell executable in your `task.yml` (see example)](https://github.com/EngineerBetter/concourse-tasks/blob/main/git-commit-if-changed/task.yml#L13-L33)
1. Please provide an [ironbird test spec](https://github.com/EngineerBetter/ironbird) in your directory, with a name ending in `_spec.yml`. PRs featuring untested tasks will not be merged.
