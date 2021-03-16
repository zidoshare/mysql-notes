import org.junit.jupiter.api.Test;
import site.zido.mysql.helpers.DBTablePrinter;

import java.sql.*;

public class MysqlTest {
    @Test
    public void testMysql() throws ClassNotFoundException, SQLException {
        Class.forName("com.mysql.jdbc.Driver");
        Connection connection_1 = DriverManager.getConnection("jdbc:mysql://mysql:3306/", "root", "123456");
        Statement statement = connection_1.createStatement();
        statement.execute("use test");
        System.out.println("1.");
        ResultSet result = statement.executeQuery("select * from t");
        DBTablePrinter.printResultSet(result);
    }
}
