package com.example.cmd;

import com.google.devtools.common.options.Option;
import com.google.devtools.common.options.OptionsBase;

public class Options extends OptionsBase {

    @Option(
            name = "address",
            abbrev = 'a',
            help = "The service server address to connect to, default address is 127.0.0.1",
            defaultValue = "127.0.0.1"
    )
    public String address;

    @Option(
            name = "port",
            abbrev = 'p',
            help = "The service port to connect to, default port is 8080",
            defaultValue = "8080"
    )
    public int port;

    @Option(
            name = "command",
            abbrev = 'c',
            help = "The command to interact with server. Valid command is list and save",
            defaultValue = "list"
    )
    public String command;

    @Option(
            name = "filter",
            abbrev = 'f',
            help = "The filter work together with command list to filter the result. Filter is format as JSON string",
            defaultValue = ""
    )
    public String filter;

    @Option(
            name = "data",
            abbrev = 'd',
            help = "The data is data to send to server with command save. Data is format as JSON string",
            defaultValue = ""
    )
    public String data;

}
