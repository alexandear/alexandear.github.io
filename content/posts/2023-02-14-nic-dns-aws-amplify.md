---
title: How to set up a free domain from NIC.UA on AWS Amplify
date: 2023-02-14
tags: ["free", "dns", "aws", "amplify"]
---

[NIC.UA](https://nic.ua/en) is a Ukrainian domain registrar that provides ".pp.ua" for [free](https://nic.ua/en/domains/.pp.ua). This domain is often used by individuals or small organizations looking for a web presence but who don't want to pay for a custom domain name.

![Free .pp.ua](/img/2023-02-14-nic-dns-aws-amplify/nic-ua-free-pp-ua.webp)

<!--more-->

[AWS Amplify](https://aws.amazon.com/amplify/) is a set of development tools and services provided by Amazon Web Services that enables developers to build and deploy web and mobile applications quickly and easily.
Unfortunately, NIC.UA doesn't have [an instruction](https://support.nic.ua/en-us/section/28-name-servers-and-dns-records) on how to configure DNS records for AWS Amplify. So, I wrote my own. Follow these steps to set up your domain.

### Configure AWS Amplify

1. Sign in to the AWS Management Console and open the [Amplify console](https://console.aws.amazon.com/amplify/).
2. For your app in the navigation pane, choose **App Settings - Domain management - Add domain**.
3. Enter the name of your root domain, and then choose **Configure domain**. E.g., `alexandear.pp.ua`.

    ![AWS add domain](/img/2023-02-14-nic-dns-aws-amplify/aws-amplify-add-domain.webp)

4. On the **Actions** menu, choose **View DNS records**. Copy all values and open the NIC.UA dashboard. Do not close the tab with the AWS Amplify console.

    ![AWS copy DNS](/img/2023-02-14-nic-dns-aws-amplify/aws-amplify-copy-dns.webp)

### Configure NIC.UA

1. *Change name servers in the domain.* Go to the order properties in the **Domains** section of your personal account. Then change **NS-servers** to the **NIC.UA name servers** item and press the **Change NS** button.

    ![NIC domain](/img/2023-02-14-nic-dns-aws-amplify/nic-ua-domain.webp)

2. *Configure DNS records on name servers.* Go to the **Name Servers** section and click on the gear-shaped button next to the renew button. Click the **Edit** button next to the **DNS Records** heading and delete all existing records. Using the **Add Record** button, create three records with parameters from AWS Amplify.

    _record 1:_
    * Name: *`_30a55502b8f33b78ce5e8f3d54d5dc36.alexandear`*
    * Type: *`CNAME`*
    * Value: *`_51bfde216354adfc5c8a9a29e8b93b7.htgdxnmnnj.acm-validations.aws.`*

    _record 2:_
    * Name: *`@`*
    * Type: *`CNAME`*
    * Value: *`d1nh7kxfyh9s3p.cloudfront.net.`*

    _record 3:_
    * Name: *`www`*
    * Type: *`CNAME`*
    * Value: *`d1nh7kxfyh9s3p.cloudfront.net.`*

    ![NIC create records](/img/2023-02-14-nic-dns-aws-amplify/nic-ua-create-dns-records.webp)

    Note that you should add the trailing dot `.` for values.

3. Wait approximately a few hours for the information to propagate to your ISP's DNS cache.
4. When AWS Amplify shows **Available** for **Status**, the domain is successfully configured and ready to use.

![AWS domain available](/img/2023-02-14-nic-dns-aws-amplify/aws-amplify-domain-status-available.webp)

Ta-da!
