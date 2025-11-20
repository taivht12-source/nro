package boss.BossManager;

/**
 *
 * @author Nguyễn tấn tài
 */
public class AnTromManager extends BossManager {

    private static AnTromManager instance;

    public static AnTromManager gI() {
        if (instance == null) {
            instance = new AnTromManager();
        }
        return instance;
    }

}