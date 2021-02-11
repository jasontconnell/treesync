*** treesync

One day I found myself working with multiple similar directory structures of config files, making a change in one then having to propagate that change to all of the other directories that were side by side.

So after doing that one time, I wrote treesync.

`treesync.json` needs to be in the root of the directory structure and just have a list of `rootFolders` that are in consideration for syncing.

```
{
    "rootFolders": ["test1", "test2", "test3", "untouched"],
    "alwaysExclude": ["untouched"],
    "log": "stdout",
    "folderGroups": {
        "group1": ["test1", "test3"]
    }
}
```

Parameters:
 - a :  The action. Supported actions are copy, delete, sync
 - exc : Exclude folders during this operation (csv)
 - inc : Include only these folders during this operation (csv)
 - file or folder : the file or folder to sync

Config File:
  - rootFolders: The folders to sync between
  - alwaysExclude: A folder that will never get files copied into it
  - folderGroups: Use this to set the include or exclude quickly without having to specify a whole list. alwaysExclude is always excluded even if included in a group.
  - log: Write logs out to a file or stdout. If a file, it will be the treesync root + the filename supplied in `log`

For delete operation, the file doesn't have to exist.

`copy` will run on either a file or a directory

`sync` will only run on a directory and will make the destination directories match file by file (it will delete all first then copy)

`delete` will run on either a file or a directory

The command can be run from anywhere (best is to go install it) and will determine the root up the file tree, much like git works with .git folder (and numerous other examples)

Example:

In a folder structure like this:
<pre>
    /test1
        /file.txt
    /test2
        /folder
            /infolder.txt
        /file2.txt
    treesync.json
</pre>

If your current working directory was `/test2/folder/` and you ran `treesync -a copy infolder.txt`, it will create the folder structure of `/test1/folder/` and copy `infolder.txt` into `/test1/folder/`.  It will do this for all folders specified in `treesync.json`.

As a bonus, you can modify the registry files included in the repo to add context menus for treesync. They will operate with what is included in the config file with no way to specify include or exclude.
