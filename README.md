# 416-team

## 0. Test Server
URL: http://120.77.220.71:8100/  

The following is list of html,css,js,image and so on.  

<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/0.png?raw=true" width="50%" height="50%" div align="center" />

The following is the main test HTML.  

URL: http://120.77.220.71:8100/c.html  

<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/3.png?raw=true" width="50%" height="50%" />
  
  
The HTML is supported by bootstrap(CDN), jquery(CDN), local Javascript and embed Javascript. When you click the "Next" or "Last" button, the javascript works and and the image src will be changed. My intention is that we can click the button and load more images into webcache and memory(of disk) size increases at the same time. If the size is greater than cache_size, then eviction policy starts to work.

## 1. Build and Run
```
go run web.go
```
The Monitoring port is 8080. We now can test the code locally and then transter it into AZURE VM.

## 2. Firefox proxy
#### Step 1 
   
<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/1.png?raw=true" width="50%" height="50%" />  

#### Step 2  
   
<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/2.png?raw=true" width="50%" height="50%" />

## 3. Run the WebCache

```
go run web.go
```
   
<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/4.png?raw=true" width="50%" height="50%" />

   
   Then open the firefox and input http://120.77.220.71:8100/c.html.
   
   <img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/3.png?raw=true" width="50%" height="50%" />
   
   <img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/5.png?raw=true" width="50%" height="50%" />
   
   The terminal will show the loading information. If the outsides resource is not in webcache, it will download it into disk and replace the old link in HTML with new ones which is linked to webcache.
   
## 4. I hope you could do
  
Now when firefox requsts a html, webcache reveive the request and do as the following:
  1. Analysis the url path of HTTP request and get the result of hashing-sha256
  2. Check whether the HTML has been in webcache(./static) now. (To do: check in the data structure in the memory)
  3. If so, read the HTML from the disk and return. (To do: return directly from HTML in the data structure in the memory)
  4. If not, download the HTML into disk and change its name by method of hash function sha256. (To do: store one copy of HTML as string in memory)
  5. Parsing the HTML, if there are outside links such as src="..." or href="...", download the source into disk and rename them in the same way.(To do: store one copy of outside source in memory such as CSS,JS,image)

<img src="https://github.com/BestOreo/Pic-for-README.md/blob/master/DistributionSystem/6.png?raw=true" width="50%" height="50%" />
