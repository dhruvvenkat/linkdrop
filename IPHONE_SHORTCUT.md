# README: Sending URLs to Your Laptop via Apple Shortcuts

In order to use this product, you need to set up an Apple Shortcut, which will make the HTTP request to our server (hosted on your PC) for you.

---

## prerequisites
* **Apple Device**: An iPhone or iPad with the native **Shortcuts** app installed.
* **Network**: Both your mobile device and your laptop **must** be connected to the same Wi-Fi network.
* **Server Details**: You will need your laptop's local IP address and the port number your server is listening on (e.g., `http://192.xxx.xxx.xxx:4545` OR `http://10.0.0.xxx:4545`).

---

## config steps

### 1. enable share sheet
This allows you to trigger the shortcut directly from Safari or any other browser.
1. Open the **Shortcuts** app on your iOS device and create a new shortcut.
2. Name your shortcut something clear, like *linkdrop* or *send URL to laptop*.
3. Tap the **Shortcut Details** icon (the "i" button or the three dots in the top right).
4. Turn on **Show in Share Sheet**.
5. Set the shortcut to receive **URLs** or **Web Pages** from the Share Sheet. This action creates a dynamic variable called `Shortcut Input`.

### 2. configure network request
1. Search for and add the **Get Contents of URL** action to your workflow.
2. In the top URL field of the action, type your laptop's server address (e.g., `http://[Your-Laptop-IP]:[Port]`).
3. Tap the **arrow icon** (or **Show More**) inside the action block to reveal advanced settings.
4. Change the **Method** from `GET` to **POST**.
5. Set the **Request Body** to **JSON**.

### 3. build JSON payload
1. Tap **Add new field** and select **Text**.
2. In the **Key** field (left side), type exactly: `url`
3. In the **Text** field (right side), tap the empty space.
4. Select **Shortcut Input** from the variable bar located right above your keyboard.

Your final action block structure will look like this:
> **Get contents of** `http://YOUR-LAPTOP-IP:PORT`
> * **Method**: POST
> * **Request Body**: JSON
> * **url** ➡️ `Shortcut Input`

---

## usage

1. Open Safari or your preferred browser on your phone.
2. Navigate to any webpage you want to share.
3. Tap the browser's **Share** button.
4. Scroll down the share menu and tap your shortcut name (*Send to Laptop*).
5. The shortcut will automatically package the current link into the following format and send it to your laptop:

```json
{
  "url": "https://example.com"
}
```
