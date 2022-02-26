// +build !ci

// +build ios

#import <UIKit/UIKit.h>

void popupNetwork() {
    [[NSURLSession.sharedSession dataTaskWithURL:[NSURL URLWithString:@"http://captive.apple.com"]] resume];
}