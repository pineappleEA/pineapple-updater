/*
Downloader for the Windows Version of pineappleEA
Copyright (C) 2020  MCredstoner2004

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA
*/


#include <iostream>
#include <string>
#include <vector>
#include <unordered_map>
#include <regex>
#include <aria2/aria2.h>
#include <curl/curl.h>
#define PINEAPPLESITE "https://raw.githubusercontent.com/pineappleEA/pineappleEA.github.io/master/index.html"
#define PINEAPPLELOGO "           /$$           /$$$$$$$$  /$$$$$$                      /$$          \n          |__/          | $$_____/ /$$__  $$                    | $$          \n  /$$$$$$  /$$ /$$$$$$$ | $$      | $$  \\ $$  /$$$$$$   /$$$$$$ | $$  /$$$$$$ \n /$$__  $$| $$| $$__  $$| $$$$$   | $$$$$$$$ /$$__  $$ /$$__  $$| $$ /$$__  $$\n| $$  \\ $$| $$| $$  \\ $$| $$__/   | $$__  $$| $$  \\ $$| $$  \\ $$| $$| $$$$$$$$\n| $$  | $$| $$| $$  | $$| $$      | $$  | $$| $$  | $$| $$  | $$| $$| $$_____/\n| $$$$$$$/| $$| $$  | $$| $$$$$$$$| $$  | $$| $$$$$$$/| $$$$$$$/| $$|  $$$$$$$\n| $$____/ |__/|__/  |__/|________/|__/  |__/| $$____/ | $$____/ |__/ \\_______/\n| $$                                        | $$      | $$                    \n| $$                                        | $$      | $$                    \n|__/                                        |__/      |__/                 \non pizza\nBrought to you by EmuWorld!\n" 

int downloadEventCallback(aria2::Session* session, aria2::DownloadEvent event, aria2::A2Gid gid, void* userData) {
	switch(event) {
	case aria2::EVENT_ON_DOWNLOAD_COMPLETE:
		std::cout << "\nDOWNLOAD COMPLETE!!!" << std::endl;
		break;
	case aria2::EVENT_ON_DOWNLOAD_ERROR:
		std::cout << "error Downloadin";
		break;
	default:
		return 0;
	}
	return 0;
}

size_t writeCallback(char* buf, size_t size, size_t count, void* userData) {
	static_cast<std::string*>(userData)->append(buf, size*count);
	return size*count;
}

struct Version{
	static const std::string linkStart;
	static const std::string linkEnd;
	static const std::string nameEnd;
	static const std::string anonStart;
	static const std::string anonEnd;
	static const std::string gDriveStart;
	static const std::string gDriveEnd;
	static const std::regex findNumber;
	std::string link;
	std::string name;
	std::string anonId;
	std::string gDriveId;
	int number;
	Version(const std::string& entry) {
		size_t start = entry.find(linkStart) + linkStart.length();
		size_t end = entry.find(linkEnd, start);
		link = entry.substr(start, end - start);
		start = end + linkEnd.length();
		end = entry.find(nameEnd, start);
		name = entry.substr(start, end - start);
		start = entry.find(gDriveStart, end);
		if(start != std::string::npos) {
		start += gDriveStart.length();
		end = entry.find(gDriveEnd, start);
		gDriveId = entry.substr(start, end - start);
		}else{
		gDriveId = "";
		}
		start = link.find(anonStart) + anonStart.length();
		end = link.find(anonEnd, start);
		anonId = link.substr(start, end - start);
		std::smatch matches;
		if(std::regex_search(name, matches, findNumber)) {
			number = std::stoi(matches[0]);
		}else{
			number = -1;
		}
		
	}
};
const std::string Version::linkStart = "<a href=";
const std::string Version::linkEnd = ">";
const std::string Version::nameEnd = "</a>";
const std::string Version::anonStart = "https://anonfiles.com/";
const std::string Version::anonEnd = "/";
const std::string Version::gDriveStart = "<!--";
const std::string Version::gDriveEnd = "-->";
const std::regex Version::findNumber = std::regex("\\d+");

int main(int argc, char** argv) {
	std::cout << PINEAPPLELOGO;
	std::vector<Version> versions;
	std::unordered_map<int,const Version*> versionsByNumber;
	curl_global_init(CURL_GLOBAL_ALL);
	CURL *curl = curl_easy_init();
	std::string pineappleData;
	curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10);
	curl_easy_setopt(curl, CURLOPT_URL, PINEAPPLESITE);
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, &writeCallback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, &pineappleData);
	int curlReturnValue = curl_easy_perform(curl);
	if(curlReturnValue != 0) {
		std::cerr << "Failed to access pineapple site, try again in a few minutes" << std::endl;
		return curlReturnValue;
	}
	size_t versionsStart = pineappleData.find("<!--link-goes-here-->");
	versionsStart = pineappleData.find('\n', versionsStart)+1;
	size_t versionsEnd = pineappleData.find("div", versionsStart);
	versionsEnd = pineappleData.rfind('\n', versionsEnd)+1;
	std::string versionData = pineappleData.substr(versionsStart, versionsEnd-versionsStart);
	pineappleData.clear();
	size_t currentPos = 0;
	size_t nextPos = 0;
	while(currentPos < versionData.length()) {
	nextPos = versionData.find('\n', currentPos);
	versions.emplace_back(versionData.substr(currentPos, nextPos - currentPos));
	currentPos = nextPos + 1;
	}
	versionData.clear();
	std::cout << versions.size() << " Versions found" << std::endl;
	std::cout << "Latest version is " << versions.front().name << std::endl;
	std::cout << " [1] Download it" << std::endl << " [2] Download another version" << std::endl << "or anything else to exit" << std::endl;
	char input;
	const Version* selectedVersion;
	std::string anonData;
	std::string anonLink;
	size_t linkStart;
	size_t linkEnd;
	std::regex findAnonLink = std::regex("https://cdn-.*?\\.anonfiles.com/[^\"]+");
	std::smatch linkMatches;
	std::cin >> input;
	long versionNumber = 0;
	std::unordered_map<int,const Version*>::const_iterator versionIterator;
	switch(input) {
		case '1':
			selectedVersion = &versions.front();
		break;
		case '2':
			std::cout << "Available versions are :" << std::endl;
			for(const Version& version : versions) {
				versionsByNumber[version.number] = &(version);
				std::cout << version.number << ", ";
			}
			std::cout << "\b\b  \b\b" << std::endl;
			std::cout << "Enter a version number" << std::endl;
			std::cin >> versionNumber;
			versionIterator = versionsByNumber.find(versionNumber);
			if(versionIterator == versionsByNumber.end()) {
				std::cerr << "version not found" << std::endl;
				return 0;
			}else {
				selectedVersion = (*versionIterator).second;
			}
		break;
		default:
		curl_easy_cleanup(curl);
		curl_global_cleanup();
		return 0;
	}
	curl_easy_setopt(curl, CURLOPT_URL, selectedVersion->link.c_str());
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, &anonData);
	curlReturnValue = curl_easy_perform(curl);
	curl_easy_cleanup(curl);
	curl_global_cleanup();
	if(curlReturnValue != 0) {
		std::cerr << "Couldn't retrieve data from anonfiles" << std::endl << "If you're on Italy or Iran, try using a VPN in another country," << std::endl;
		std::cerr << "otherwise, please try again in a few minutes" << std::endl;
		return curlReturnValue;
	}
	if(!std::regex_search(anonData, linkMatches, findAnonLink)) {
		std::cerr << "can't find link inside anonfiles" << std::endl;
		return 1;
	}
	anonLink = linkMatches[0];
	aria2::libraryInit();
	aria2::Session* session;
	aria2::SessionConfig config;
	config.downloadEventCallback = downloadEventCallback;
	aria2::KeyVals options;
	options.push_back(std::pair<std::string, std::string> ("max-connection-per-server", "6"));
	options.push_back(std::pair<std::string, std::string> ("split", "12"));
	options.push_back(std::pair<std::string, std::string> ("continue", "true"));
	session = aria2::sessionNew(options, config);
	std::vector<std::string> uris = {anonLink};
	int rv = aria2::addUri(session, nullptr, uris, options);
	std::string outputLine;
	size_t lastLength = 0;
	do{
	rv = aria2::run(session, aria2::RUN_ONCE);
		for(aria2::A2Gid gid : aria2::getActiveDownload(session)) {
			aria2::DownloadHandle* dh = aria2::getDownloadHandle(session, gid);
			if(dh) {	
				int64_t current = dh->getCompletedLength();
				int64_t total = dh->getTotalLength();
				if(total > 0) {
					int64_t progress = 51 * current / total;
					outputLine = "Downloading [";
					if(progress > 0) {
						outputLine.insert(outputLine.length(),progress-1, '-');
						outputLine += ">";
					}
					if(progress < 51) {
						outputLine.insert(outputLine.length(), 50 - progress, '.');
					}
					outputLine += "] " + std::to_string(100 * current / total) + "%";
					outputLine += " (" + std::to_string(dh->getDownloadSpeed()/1024) + "KiB/s)";
					std::cout << std::string(lastLength, '\b');
					std::cout << std::string(lastLength, ' ');
					std::cout << std::string(lastLength, '\b');
					std::cout << outputLine << std::flush;
					lastLength = outputLine.length();					
				}
				aria2::deleteDownloadHandle(dh);
			}
		}
	}while(rv == 1);
	rv = aria2::sessionFinal(session);
	aria2::libraryDeinit();
	return rv;
}
