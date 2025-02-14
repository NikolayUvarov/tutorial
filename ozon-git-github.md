

# GIT usage in github
## Inital steps
Create directory and do inside

    git init 
    git add . 
    git commit -m comment
    git commit -m "comment with multiwords"

Note that git makes default branch **master**, but github for some vague reason uses **main**.
So create branch main locally and go to it and remove master:
git branch main
git checkout main
git branch -d master






## Use remote repo, e.g. github
Add remote repository

    git remote add origin  git@github.com:NikolayUvarov/tutorial.git

Check remote 

    git remote -v

Set remote branch if using defalut master branch locally

    git branch --set-upstream-to=origin/main master

Set remote branch if using main to repo look like in github







### tu
#### jj

hello


