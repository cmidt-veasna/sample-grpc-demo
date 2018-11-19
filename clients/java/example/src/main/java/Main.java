import com.example.cmd.List;
import com.example.cmd.Options;
import com.example.cmd.Save;
import com.google.devtools.common.options.OptionsParser;

import java.util.Collections;

public class Main {

    public static void main(String[] args) {
        OptionsParser parser = OptionsParser.newOptionsParser(Options.class);
        parser.parseAndExitUponError(args);
        Options options = parser.getOptions(Options.class);

        if (options.command.isEmpty() && options.data.isEmpty()) {
            parser.describeOptions(Collections.emptyMap(), OptionsParser.HelpVerbosity.LONG);
            return;
        }

        switch (options.command) {
            case "list":
                new List(options).list();
                break;

            case "save":
                new Save(options).save();
                break;
        }
    }

}
