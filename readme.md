*** treesync

One day I found myself working with multiple similar directory structures of config files, making a change in one then having to propagate that change to all of the other directories that were side by side.

So after doing that one time, I wrote treesync.

treesync.json needs to be in the root of the directory structure and just have a list of "rootFolders" that are in consideration for syncing.

```
{
    "rootFolders": ["test1", "test2"],
    "alwaysExclude": []
}
```

Parameters:
 - a :  The action. Supported actions are copy, delete, sync
 - e : Exclude folder during this operation
 - file or folder : the file or folder to sync

For delete operation, the file doesn't have to exist.

`copy` will run on either a file or a directory

`sync` will only run on a directory and will make the destination directories match file by file (it will delete all first then copy)

`delete` will only work on files currently

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