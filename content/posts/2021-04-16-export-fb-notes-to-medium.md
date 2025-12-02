---
title: How to export Facebook Notes to Medium with the original publication date
date: 2021-04-16
tags: ["facebook", "medium"]
---

Recently I was getting to know [Facebook discontinued its Notes](https://www.digitalinformationworld.com/2020/10/facebook-is-quietly-sunsetting-its-notes-feature.html) feature for Individuals and Pages.
I have a few Facebook notes and need to move them to another live blog platform, e.g. Medium.
But there is no way to import from Facebook with the help of [Medium's import functionality](https://help.medium.com/hc/en-us/articles/214550207-Import-a-post).
The only possible solution is to manually create a story, copy a note's content, and set a date to the original publication date.

The trickier step is keeping the original publication date. Below is how I propose to do this.

{{< figure src="/img/2021-04-16-export-fb-notes-to-medium/imported-fb-note.webp" width="100%" alt="A Facebook note imported to a Medium story with the preserved publish date" >}}

<!--more-->

{{< note >}}
*Originally published on [Medium](https://alexandear.medium.com/%D0%BC%D1%96%D0%B9-%D0%BF%D0%B5%D1%80%D1%88%D0%B8%D0%B9-%D1%80%D0%B0%D0%B7-9992b4a49b8f).*
{{< /note >}}

## Find Facebook notes

First, let's find our Facebook notes to export. Go to profile — three dots — "Activity log" — "Filter" — choose the category "Notes". Or open the link, replacing `YOUR_FB_USERNAME`:

```txt
https://www.facebook.com/YOUR_FB_USERNAME/allactivity/?category_key=NOTECLUSTER&filter_hidden=ALL&filter_privacy=NONE&manage_mode=false
```

{{< figure src="/img/2021-04-16-export-fb-notes-to-medium/facebook-note.webp" width="100%" caption="Facebook activity log with Notes filter" >}}

Next, choose a note and copy its content with Control-C.

## Create a temporary WordPress post

Unlike Medium, WordPress allows setting a publication date. Go to WordPress home — "Posts" — "Add New" — past previously copied note's content. Obviously, you had to have a WordPress account.

Next, open "Settings" — "Post" — set appropriate "Publish" date.

{{< figure src="/img/2021-04-16-export-fb-notes-to-medium/wordpress-add-content.webp" width="100%" caption="Add content and set publish date for a WordPress blog" >}}

Click "Publish" and make sure the post "Visibility" is "Public". Copy the published post URL.

## Import the WordPress post to Medium

Open your Medium homepage — "Stories" — ["Import a story"](https://help.medium.com/hc/en-us/articles/214550207-Importing-a-post-to-Medium) — paste WordPress's copied URL — click "Import".

{{< figure src="/img/2021-04-16-export-fb-notes-to-medium/medium-import-wordpress.webp" width="100%" caption="Import the WordPress post to Medium" >}}

Finally, edit your story, add tags, and click "Publish".

{{< figure src="/img/2021-04-16-export-fb-notes-to-medium/imported-fb-note.webp" width="100%" caption="A Facebook note imported to a Medium story with the preserved publish date" >}}

That it. Now [the Facebook note](https://www.facebook.com/notes/455167882120716/) imported to [Medium](https://alexandear.medium.com/export-fb-notes-medium-original-date-adb21a4ad8b#:~:text=note%20imported%20to-,Medium,-with%20the%20original) with the original publish date.
