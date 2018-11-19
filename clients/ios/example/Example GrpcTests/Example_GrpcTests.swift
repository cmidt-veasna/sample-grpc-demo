//
//  Example_GrpcTests.swift
//  Example GrpcTests
//
//  Created by Veasna Sreng on 10/26/18.
//  Copyright © 2018 Veasna Sreng. All rights reserved.
//

import XCTest
@testable import Example_Grpc

class Example_GrpcTests: XCTestCase {

    override func setUp() {
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }

    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testGrpcExample() {
        // This is an example of a functional test case.
        // Use XCTAssert and related functions to verify your tests produce the correct results.

        let service = Example_ElementServiceServiceClient(address: "127.0.0.1:8080", secure: false)
        
        // list element
        var filter = Example_ElementFilter()
        filter.age = "[10,34]"
        do {
            let ele = try? service.listElement(filter)
            XCTAssert(ele != nil)
            XCTAssert(ele?.elements != nil)
        } catch {
        }
        
        // save element
        var elem = Example_Element()
        elem.name = "Test Element 1"
        elem.age = 32
        elem.status = 8
        do {
            let ele = try? service.persistElement(elem)
            XCTAssert(ele != nil)
            XCTAssert(ele?.id != "")
            XCTAssert(ele?.createdAt != "")
            XCTAssert(ele?.updatedAt != "")
        } catch{
        }
    }

    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }

}
